import { Colors, Typography } from '../../styles'
import { Controller, useForm } from 'react-hook-form'
import React, { useState } from 'react'

import Cookies from 'js-cookie'
import GoogleSignInButton from '../atoms/buttons/GoogleSignInButton'
import JoinWaitlistButton from '../atoms/buttons/JoinWaitlistButton'
import { Navigate } from 'react-router-dom'
import UnauthorizedFooter from '../molecules/UnauthorizedFooter'
import UnauthorizedHeader from '../molecules/UnauthorizedHeader'
import apiClient from '../../utils/api'
import styled from 'styled-components'
import { useAppSelector } from '../../redux/hooks'
import { AUTHORIZATION_COOKE } from '../../constants'

const LandingScreenContainer = styled.div`
    background-color: ${Colors.white};
    height: 100vh;
    display: flex;
    flex-direction: column;
    font-family: Switzer-Variable;
`
const FlexColumn = styled.div`
    display: flex;
    flex-direction: column;
`
const FlexGrowContainer = styled.div`
    flex: 1;
`
const Header = styled.div`
    max-width: 700px;
    margin: auto;
    margin-bottom: 40px;
    font-size: ${Typography.landingScreen.header};
    text-align: center;
    font-family: inherit;
`
const Subheader = styled.div`
    max-width: 725px;
    margin: auto;
    font-size: ${Typography.landingScreen.subheader};
    text-align: center;
`
const WaitlistContainer = styled.div`
    display: flex;
    flex-direction: row;
    height: 35px;
    margin: auto;
    margin-top: 30px;
    width: 500px;
`
const WaitlistTextInput = styled.input`
    outline: none;
    flex: 1;
    padding: 0px 10px;
`
const ResponseContainer = styled.div`
    display: flex;
    justify-content: center;
    margin-top: 20px;
    height: 20px;
    color: ${Colors.response.error};
`

const LandingScreen = () => {
    const [message, setMessage] = useState('')
    const { control, handleSubmit } = useForm({
        defaultValues: {
            email: '',
        },
    })
    const onWaitlistSubmit = (data: { email: string }) => {
        joinWaitlist(data.email)
    }
    const onWaitlistError = () => {
        setMessage('Email field is required')
    }

    const joinWaitlist = async (email: string) => {
        const response: Response = await apiClient.post('/waitlist/', {
            email: email,
        })
        if (response.status === 201) {
            setMessage('Success!')
        } else if (response.status === 302) {
            setMessage('This email already exists in the waitlist')
        } else {
            setMessage('There was an error adding you to the waitlist')
        }
    }
    const { authToken } = useAppSelector((state) => ({ authToken: state.user_data.auth_token }))
    const authCookie = Cookies.get(AUTHORIZATION_COOKE)

    if (authToken || authCookie) return <Navigate to="/tasks" />

    return (
        <LandingScreenContainer>
            <UnauthorizedHeader />
            <FlexGrowContainer>
                <FlexColumn>
                    <Header>The task manager for highly productive people.</Header>
                    <Subheader>
                        General Task pulls together your emails, messages, and tasks and prioritizes what matters most.{' '}
                    </Subheader>
                    <Subheader></Subheader>
                </FlexColumn>
                <WaitlistContainer>
                    <Controller
                        control={control}
                        rules={{
                            required: true,
                        }}
                        render={({ field: { onChange, value } }) => (
                            <WaitlistTextInput
                                type="text"
                                onChange={onChange}
                                value={value}
                                placeholder="Enter email address"
                            />
                        )}
                        name="email"
                    />
                    <JoinWaitlistButton onSubmit={handleSubmit(onWaitlistSubmit, onWaitlistError)} />
                </WaitlistContainer>
                <ResponseContainer data-testid="response-container">{message}</ResponseContainer>
                <GoogleSignInButton />
            </FlexGrowContainer>
            <UnauthorizedFooter />
        </LandingScreenContainer>
    )
}

export default LandingScreen
import React, { useEffect } from 'react'
import { useQuery } from 'react-query'
import { Navigate } from 'react-router-dom'
import Loading from '../components/atoms/Loading'
import SingleViewTemplate from '../components/templates/SingleViewTemplate'
import ModalView from '../components/views/ModalView'
import TermsOfServiceView from '../components/views/TermsOfServiceSummaryView'
import { useAppDispatch } from '../redux/hooks'
import { setShowModal } from '../redux/tasksPageSlice'
import { fetchUserInfo } from '../services/queryUtils'
import { ModalEnum } from '../utils/enums'

const TermsOfServiceSummaryScreen = () => {
    const dispatch = useAppDispatch()
    const { data, isLoading } = useQuery('user_info', fetchUserInfo)
    useEffect(() => {
        dispatch(setShowModal(ModalEnum.PRIVACY_POLICY))
    }, [])

    if (isLoading) return <Loading />
    if (!isLoading && data.agreed_to_terms) return <Navigate to="/" />
    return (
        <SingleViewTemplate>
            <ModalView size="medium" canClose={false}>
                <TermsOfServiceView />
            </ModalView>
        </SingleViewTemplate>
    )
}

export default TermsOfServiceSummaryScreen
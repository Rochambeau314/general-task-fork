import { useEffect, useLayoutEffect, useRef } from 'react'
import styled from 'styled-components'
import { Border, Spacing, Typography } from '../../../styles'
import { FontSize, PlainTextEditorProps } from './types'

const PlainTextArea = styled.textarea<{ fontSize: FontSize }>`
    background-color: inherit;
    outline: none;
    border: none;
    border-radius: ${Border.radius.small};
    resize: none;
    width: 100%;
    box-sizing: border-box;
    height: 100%;
    padding: ${Spacing._8};
    white-space: pre-wrap;
    ${({ fontSize }) => fontSize === 'small' && Typography.bodySmall};
    ${({ fontSize }) => fontSize === 'medium' && Typography.subtitle};
    ${({ fontSize }) => fontSize === 'large' && Typography.title};
`

const PlainTextEditor = (props: PlainTextEditorProps) => {
    const { isFullHeight, maxHeight, value, onChange, ...rest } = props
    const ref = useRef<HTMLTextAreaElement>(null)
    const resizeEditor = () => {
        if (!isFullHeight && ref.current) {
            ref.current.style.height = '0px'
            ref.current.style.height =
                maxHeight && ref.current.scrollHeight > maxHeight ? `${maxHeight}px` : `${ref.current.scrollHeight}px`
        }
    }

    useLayoutEffect(() => {
        if (ref.current) {
            ref.current.value = value
        }
    }, [value])

    useEffect(resizeEditor, [value, maxHeight])

    useEffect(() => {
        if (props.autoSelect && ref.current) {
            ref.current.select()
        }
    }, [])

    const handleKeyDown: React.KeyboardEventHandler<HTMLTextAreaElement> = (e) => {
        if (ref.current && (e.key === 'Escape' || (props.blurOnEnter && e.key === 'Enter'))) {
            ref.current.blur()
        }
    }

    return (
        <PlainTextArea
            ref={ref}
            onChange={(e) => {
                resizeEditor()
                onChange(e.target.value)
            }}
            onKeyDown={handleKeyDown}
            {...rest}
        />
    )
}

export default PlainTextEditor
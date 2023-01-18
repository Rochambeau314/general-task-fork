import { useState } from 'react'
import { DateTime } from 'luxon'
import styled from 'styled-components'
import { Colors, Spacing, Typography } from '../../styles'
import { icons } from '../../styles/images'
import { TNote } from '../../utils/types'
import { getHumanDateTime } from '../../utils/utils'
import Flex from '../atoms/Flex'
import { Icon } from '../atoms/Icon'
import TaskTemplate from '../atoms/TaskTemplate'
import { Label, Truncated } from '../atoms/typography/Typography'
import ItemContainer from '../molecules/ItemContainer'
import NoteContextMenuWrapper from './NoteContextMenuWrapper'

const NoteTitle = styled(Truncated)`
    ${Typography.bodySmall};
`
const TitleContainer = styled.span<{ deleted?: boolean }>`
    display: flex;
    gap: ${Spacing._8};
    align-items: center;
    min-width: 0;
    margin-right: ${Spacing._8};
    text-decoration: ${({ deleted }) => (deleted ? 'line-through' : 'none')};
    color: ${({ deleted }) => (deleted ? Colors.text.light : Colors.text.black)};
`
interface NoteProps {
    note: TNote
    isSelected: boolean
    onSelect: (note: TNote) => void
}
const Note = ({ note, isSelected, onSelect }: NoteProps) => {
    const [contextMenuOpen, setContextMenuOpen] = useState(false)
    const isShared = +DateTime.fromISO(note.shared_until ?? '0') > +DateTime.local()
    return (
        <NoteContextMenuWrapper note={note} onOpenChange={setContextMenuOpen}>
            <TaskTemplate>
                <ItemContainer isSelected={isSelected} onClick={() => onSelect(note)} forceHoverStyle={contextMenuOpen}>
                    <TitleContainer deleted={note.is_deleted}>
                        <Icon icon={icons.note} />
                        <NoteTitle>{note.title}</NoteTitle>
                    </TitleContainer>
                    <Flex gap={Spacing._12} alignItems="center">
                        {isShared && <Icon icon={icons.link} />}
                        <Label color="light">{getHumanDateTime(DateTime.fromISO(note.created_at))}</Label>
                    </Flex>
                </ItemContainer>
            </TaskTemplate>
        </NoteContextMenuWrapper>
    )
}

export default Note

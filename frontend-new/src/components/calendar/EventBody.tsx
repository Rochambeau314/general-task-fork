import {
    EventBodyStyle,
    EventInfoContainer,
    EventInfo,
    EventTitle,
    EventTime,
    EventFillContinues,
    EventFill,
    CELL_HEIGHT,
} from './CalendarEvents-styles'
import React from 'react'
import { TEvent } from '../../utils/types'

const LONG_EVENT_THRESHOLD = 45 // minutes
const MINIMUM_BODY_HEIGHT = 15 // minutes

interface EventBodyProps {
    event: TEvent
    collisionGroupSize: number
    leftOffset: number
}
function EventBody({ event, collisionGroupSize, leftOffset }: EventBodyProps): JSX.Element {
    const startTime = new Date(event.datetime_start)
    const endTime = new Date(event.datetime_end)
    const timeDurationMinutes = (endTime.getTime() - startTime.getTime()) / 1000 / 60
    const timeDurationHours = timeDurationMinutes / 60

    const rollsOverMidnight = endTime.getDay() !== startTime.getDay()
    const eventBodyHeight = Math.max(
        MINIMUM_BODY_HEIGHT,
        rollsOverMidnight
            ? ((new Date(startTime).setHours(24, 0, 0, 0) - startTime.getTime()) / 1000 / 3600) * CELL_HEIGHT
            : timeDurationHours * CELL_HEIGHT
    )

    const startTimeHours = startTime.getHours() - 1
    const startTimeMinutes = startTime.getMinutes()
    const topOffset = (60 * startTimeHours + startTimeMinutes) * (CELL_HEIGHT / 60)

    const MMHH = { hour: 'numeric', minute: 'numeric', hour12: true } as const
    const startTimeString = startTime.toLocaleString('en-US', MMHH).replace(/AM|PM/, '')
    const endTimeString = endTime.toLocaleString('en-US', MMHH)

    const isLongEvent = timeDurationMinutes >= LONG_EVENT_THRESHOLD
    const eventHasEnded = endTime.getTime() < Date.now()

    return (
        <EventBodyStyle
            key={event.id}
            squishFactor={collisionGroupSize}
            leftOffset={leftOffset}
            topOffset={topOffset}
            eventBodyHeight={eventBodyHeight}
            eventHasEnded={eventHasEnded}
        >
            <EventInfoContainer>
                <EventInfo isLongEvent={isLongEvent}>
                    <EventTitle isLongEvent={isLongEvent}>{event.title || '(no title)'}</EventTitle>
                    <EventTime>{`${startTimeString} - ${endTimeString}`}</EventTime>
                </EventInfo>
            </EventInfoContainer>
            {rollsOverMidnight ? <EventFillContinues /> : <EventFill />}
        </EventBodyStyle>
    )
}

export default EventBody
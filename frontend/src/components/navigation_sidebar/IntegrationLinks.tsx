import { useMemo } from 'react'
import { useLocation } from 'react-router-dom'
import { usePreviewMode } from '../../hooks'
import { useGetPullRequests } from '../../services/api/pull-request.hooks'
import { useGetLinkedAccounts } from '../../services/api/settings.hooks'
import { useGetTasks } from '../../services/api/tasks.hooks'
import { icons, logos } from '../../styles/images'
import { PR_SORT_AND_FILTER_CONFIG } from '../../utils/sortAndFilter/pull-requests.config'
import useSortAndFilterSettings from '../../utils/sortAndFilter/useSortAndFilterSettings'
import { TPullRequest } from '../../utils/types'
import { doesAccountNeedRelinking, isGithubLinked, isLinearLinked, isSlackLinked } from '../../utils/utils'
import NavigationLink from './NavigationLink'

interface IntegrationLinksProps {
    isCollapsed?: boolean
}
const IntegrationLinks = ({ isCollapsed }: IntegrationLinksProps) => {
    const { data: pullRequestRepositories } = useGetPullRequests()
    const { pathname } = useLocation()
    const { data: folders } = useGetTasks()
    const { isPreviewMode } = usePreviewMode()

    const linearTasksCount = useMemo(() => {
        const tasks =
            folders?.filter((section) => !section.is_done && !section.is_trash).flatMap((folder) => folder.tasks) ?? []
        return tasks.filter((task) => task.source.name === 'Linear').length
    }, [folders])

    const slackTasksCount = useMemo(() => {
        const tasks =
            folders?.filter((section) => !section.is_done && !section.is_trash).flatMap((folder) => folder.tasks) ?? []
        return tasks.filter((task) => task.source.name === 'Slack' && (!task.is_done || task.optimisticId)).length
    }, [folders])

    const { data: linkedAccounts } = useGetLinkedAccounts()

    const isGithubIntegrationLinked = isGithubLinked(linkedAccounts || [])
    const isLinearIntegrationLinked = isLinearLinked(linkedAccounts || [])
    const isSlackIntegrationLinked = isSlackLinked(linkedAccounts || [])

    const { selectedFilter: pullRequestFilter } = useSortAndFilterSettings<TPullRequest>(PR_SORT_AND_FILTER_CONFIG)
    const githubCount = isGithubIntegrationLinked
        ? pullRequestRepositories?.reduce<number>(
              (total, repo) => total + repo.pull_requests.filter(pullRequestFilter.lambda).length,
              0
          )
        : undefined
    const linearCount = isLinearIntegrationLinked ? linearTasksCount : undefined
    const slackCount = isSlackIntegrationLinked ? slackTasksCount : undefined
    return (
        <>
            <NavigationLink
                link="/overview"
                title="Overview"
                icon={icons.list}
                isCurrentPage={pathname.split('/')[1] === 'overview'}
                isCollapsed={isCollapsed}
            />
            {isPreviewMode && (
                <NavigationLink
                    link="/recurring-tasks"
                    title="Recurring tasks"
                    icon={icons.arrows_repeat}
                    iconColor="green"
                    isCurrentPage={pathname.split('/')[1] === 'recurring-tasks'}
                    isCollapsed={isCollapsed}
                />
            )}
            <NavigationLink
                link="/focus-mode"
                title="Enter Focus Mode"
                icon={icons.headphones}
                isCurrentPage={pathname.split('/')[1] === 'focus-mode'}
                isCollapsed={isCollapsed}
            />
            <NavigationLink
                link="/pull-requests"
                title="GitHub"
                icon={logos.github}
                count={githubCount}
                needsRelinking={doesAccountNeedRelinking(linkedAccounts || [], 'GitHub')}
                isCurrentPage={pathname.split('/')[1] === 'pull-requests'}
                isCollapsed={isCollapsed}
            />
            <NavigationLink
                link="/linear"
                title="Linear"
                icon={logos.linear}
                count={linearCount}
                needsRelinking={doesAccountNeedRelinking(linkedAccounts || [], 'Linear')}
                isCurrentPage={pathname.split('/')[1] === 'linear'}
                isCollapsed={isCollapsed}
            />
            <NavigationLink
                link="/slack"
                title="Slack"
                icon={logos.slack}
                count={slackCount}
                needsRelinking={doesAccountNeedRelinking(linkedAccounts || [], 'Slack')}
                isCurrentPage={pathname.split('/')[1] === 'slack'}
                isCollapsed={isCollapsed}
            />
        </>
    )
}

export default IntegrationLinks

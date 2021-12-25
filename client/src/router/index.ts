import User from '@/models/User'
import {createRouter, createWebHashHistory, RouteLocationNormalized, RouteLocationRaw, RouteRecordRaw} from 'vue-router'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Contests',
        component: () => import("@/views/Contests.vue")
    },
    {
        path: '/about/',
        name: 'About',
        component: () => import('@/views/About.vue')
    },
    {
        path: '/profile/',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        children: [
            {
                path: 'contests/',
                name: 'contests',
                component: () => import("@/components/director/DirectorContests.vue")
            },
            {
                path: 'organizers/',
                name: 'organizers',
                component: () => import("@/components/director/DirectorOrganizers.vue")
            },
            {
                path: 'other/',
                name: 'other',
                component: () => import("@/components/director/DirectorOther.vue")
            },
            {
                path: 'org-teams/',
                name: 'org-teams',
                component: () => import("@/components/director/DirectorTeams.vue")
            },
            {
                path: 'attendance/',
                name: "attendance",
                component: () => import("@/components/user/AttendanceList.vue")
            }
        ]
    },
    {
        path: '/register-contestant/',
        name: 'ContestantSignup',
        component: () => import('@/components/contestant/ContestantSignup.vue'),
        props: (route: RouteLocationNormalized) => ({
            profile: new User(), // like the props declaration in the vue component
            ...route.params
        })
    },
    {
        path: '/notifications/',
        name: 'Notifications',
        component: () => import('@/views/Notifications.vue')
    },
    {
        path: "/contest/", // ?id=contestId
        name: "Contest",
        component: () => import("@/components/contest/ContestPage.vue"),
        children: [
            {
                name: "details",
                path: "details/", // ?id=contestId
                component: () => import("../components/contest/ContestDetails.vue")
            },
            {
                name: "teams",
                path: "teams/", // ?id=contestId
                component: () => import("../components/contest/ContestTeams.vue")
            },
            {
                name: "teamless",
                path: "teamless/", // ?id=contestId
                component: () => import("../components/contest/ContestTeamless.vue")
            },
            {
                name: "support",
                path: "support/", // ?id=contestId
                component: () => import("../components/contest/ContestSupport.vue")
            },
        ]
    },
    {
        path: '/finish-org-profile/',
        name: 'OrganizerSignup',
        component: () => import('@/components/organizer/OrganizerSignup.vue')
    },
    {
        path: '/admin/',
        name: 'Admin',
        component: () => import('@/views/Admin.vue')
    },
]

const router = createRouter({
    history: createWebHashHistory(process.env.BASE_URL),
    routes: routes,
})

export default router

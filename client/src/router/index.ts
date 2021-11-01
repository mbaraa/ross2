import {createRouter, createWebHistory, RouteRecordRaw} from 'vue-router'

const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Contests',
        component: () => import("@/views/Contests.vue")
    },
    {
        path: '/about',
        name: 'About',
        component: () => import('@/views/About.vue')
    },
    {
        path: '/profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue')
    },
    {
        path: '/notifications',
        name: 'Notifications',
        component: () => import('@/views/Notifications.vue')
    },
    {
        path: "/contest/", // ?id=contestId
        name: "Contest Details - Ross 2",
        component: () => import("@/components/ContestPage.vue"),
        children: [
            {
                name: "details",
                path: "details/", // ?id=contestId
                component: () => import("../components/ContestDetails.vue")
            },
            {
                name: "teams",
                path: "teams/", // ?id=contestId
                component: () => import("../components/ContestTeams.vue")
            },
            {
                name: "teamless",
                path: "teamless/", // ?id=contestId
                component: () => import("../components/ContestTeamless.vue")
            },
            {
                name: "support",
                path: "support/", // ?id=contestId
                component: () => import("../components/ContestSupport.vue")
            },
        ]
    },
    {
        path: '/finish-profile/',
        name: 'Finish Profile',
        component: () => import('@/components/ContestantSignup.vue')
    },
]

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes
})

export default router

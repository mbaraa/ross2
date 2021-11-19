<template>
    <div>
        &nbsp;{{ getRemainingTime() }}&nbsp;
    </div>
</template>

<script lang="ts">

import {defineComponent} from "vue";

export default defineComponent({
    props: {
        endTimestamp: Date
    },
    data() {
        return {
            name: "TimerCountdown",
            days: 0,
            hours: 0,
            minutes: 0,
            seconds: 0,
        }
    },
    computed: {
        _seconds: () => 1000,
        _minutes() {
            return this._seconds * 60
        },
        _hours() {
            return this._minutes * 60
        },
        _days() {
            return this._hours * 24
        }
    },
    methods: {
        formatNumber: function (num: number) : string {
            return num < 10 ? "0" + num.toString(10) : num.toString(10);
        },
        calcRemainingTime: function () {
            const timer = setInterval(() => {
                const now = new Date();
                const end = new Date(this.endTimestamp);
                const time = end.getTime() - now.getTime();
                if (time < 0) {
                    clearInterval(timer)
                    return;
                }
                const days = Math.floor(time / this._days);
                const hours = Math.floor((time % this._days) / this._hours);
                const minutes = Math.floor((time % this._hours) / this._minutes);
                const seconds = Math.floor((time % this._minutes) / this._seconds);
                this.days = this.formatNumber(days);
                this.hours = this.formatNumber(hours);
                this.minutes = this.formatNumber(minutes);
                this.seconds = this.formatNumber(seconds);
            }, 1000)
        },
        getRemainingTime: function () {
            return (this.days !== 0 && this.hours !== 0 && this.minutes !== 0 && this.seconds !== 0) ?
                `${this.days}:${this.hours}:${this.minutes}:${this.seconds}` :
                "OVER!";
        }
    },
    mounted() {
        this.calcRemainingTime()
    }
});
</script>

<style scoped>
</style>

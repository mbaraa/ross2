export function getLocaleTime(time: Date): string {
    return new Date(time).toLocaleTimeString("en-US", {
        hour12: true,
        hour: "2-digit",
        minute: "2-digit",
        year: "numeric",
        month: "short",
        day: "2-digit",
        weekday: "short"
    })
}

export function formatDuration(minutes: number): string {
    minutes /= 100000000000
    const hours = (minutes / 60) % 24;
    const minutes1 = minutes % 60;

    return `${hours} hours & ${minutes1} minutes`;
}

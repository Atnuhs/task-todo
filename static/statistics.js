const init = async () => {
    const container = document.getElementById("events-container")

    const today = new Date()

    const oneWeekAgo = new Date()
    oneWeekAgo.setDate(today.getDate() - 7)

    const formatDate = (date) => {
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, "0")
        const day = String(date.getDate()).padStart(2, "0")
        return `${year}-${month}-${day}`
    };

    const startDate = formatDate(oneWeekAgo);
    const endDate = formatDate(today);

    const apiUrl = `/tasks/statistics/done?start_date=${startDate}&end_date=${endDate}`;
    const resp = await fetch(apiUrl)

    if (!resp.ok) {
        console.error("failed to get events")
        return
    }

    const statisics = await resp.json();
    formattedJson = JSON.stringify(events, null, 4)
    console.log(formattedJson)
    container.textContent = formattedJson
}

document.addEventListener('DOMContentLoaded', async () => {
    init()
})

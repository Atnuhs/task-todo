const init = async () => {
    const container = document.getElementById("events-container")
    container.style.whiteSpace = "pre-wrap"
    const resp = await fetch("/events", {method: "get"})

    if (!resp.ok) {
        console.error("failed to get events")
        return
    }

    const events = await resp.json();
    formattedJson = JSON.stringify(events, null, 4)
    console.log(formattedJson)
    container.textContent = formattedJson
}

document.addEventListener('DOMContentLoaded', async () => {
    init()
})

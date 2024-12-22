const ChangeStatusButton = () => {
    const button = document.createElement("button")
    Object.assign(button.style, {
        backgroundColor: "#4CAF50",
        color: "white",
        padding: "10px 20px",
        marginRight: "1em",
        border: "none",
        borderRadius: "5px",
        cursor: "pointer",
        fontSize: "16px",
        transition: "background-color 0.3s",
    });
    return button
}

const showTask = (task) => {
    const container = document.getElementById("taskItemsArea")
    const title = document.createElement("p")
    title.textContent = task.name

    const taskItem = document.createElement("div")
    Object.assign(taskItem.style, {
        backgroundColor: "#eeeeee",
        padding: "1em",
        marginTop: "1em",
        marginBottom: "1em",
    })
    taskItem.appendChild(title)


    const status = document.createElement("p")
    status.textContent = task.status
    taskItem.appendChild(status)

    const updateButton = (btn, status) => {
        if (task.status === status) {
            btn.disabled = true
            btn.style.backgroundColor = "#ccc"
            btn.style.cursor = "default"
        } else {
            btn.disabled = false;
            btn.style.backgroundColor = "#4CAF50";
            btn.style.cursor = "pointer";
        }
    }
    
    const btnStatus = [
        {display_text: "Start", status: "doing", btn: ChangeStatusButton()},
        {display_text: "Stop", status: "pending", btn: ChangeStatusButton()},
        {display_text: "Finish", status: "completed", btn: ChangeStatusButton()},
        {display_text: "Cancel", status: "cancelled", btn: ChangeStatusButton()},
    ]
    
    btnStatus.forEach(d => {
        d.btn.textContent = d.display_text
        updateButton(d.btn, d.status)
        d.btn.addEventListener("click", async () => {
            if (d.btn.disabled) return

            const resp = await fetch(`/tasks/${task.task_id}`, {
                method: "PATCH",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({status: d.status}),
            })

            if (resp.ok) {
                task.status = d.status
                status.textContent = d.status
                btnStatus.forEach(d2 => {
                    updateButton(d2.btn, d2.status)
                })
            } else {
                alert("Failed to update task status")
            }
        })
        taskItem.appendChild(d.btn)
    });
    container.insertBefore(taskItem, container.firstChild)
}

const init = async () => {
    const form = document.getElementById("taskForm")    

    form.addEventListener("submit", async(event) => {
        event.preventDefault()

        const formData = new FormData(form)
        const taskName = formData.get("taskName")

        const resp = await fetch("/tasks", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({name: taskName})
        })

        if (!resp.ok) {
            console.error("Failed to add task")
            return
        }

        const newTask = await resp.json();
        showTask(newTask)

        form.reset()
    })

    const resp = await fetch("/tasks")
    const tasks = await resp.json()
    tasks.forEach(t => showTask(t))
}


document.addEventListener('DOMContentLoaded', async () => {
    init()
})
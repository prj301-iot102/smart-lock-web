async function listDoors() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("doorsTableBody");
    const message = document.getElementById("doorsMsg");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/door/", 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
        });
        if (!response.ok) {return;}

        const data = await response.json();

        data.forEach(door => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${door.id}</td>
                <td>${door.door_name}</td>
                <td>${door.device_id}</td>
                <td>${timeFormat(door.created_at)}</td>
                <td>${timeFormat(door.updated_at)}</td>
                <td>
                    <div class="roleDropdown">
                        <button onclick="toggleRoleDropdown(this)">Permission</button>
                        <div class="roleDropdownContent" data-door-id="${door.id}"></div>
                    </div>
                </td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error listing doors: ", error);
    }
}

async function getDoorByID(id) {
    const token = localStorage.getItem("token");

    if (!id) return null;

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/door/${id}`, {
            method: "GET",
            headers: {
                "Authorization": `Bearer ${token}`,
                "Accept": "application/json"
            }
        });

        if (!response.ok) return null;
        return await response.json();
    } catch (error) {
        console.error("Error fetching door:", error);
        return null;
    }
}

async function loadRolesForDoors() {
    const token = localStorage.getItem("token");
    if (!token) return;

    try {
        const roleRes = await fetch("https://smart-lock.patohru.qzz.io/api/role/", {
            headers: {
                "Authorization": `Bearer ${token}`,
                "Accept": "application/json"
            }
        });

        if (!roleRes.ok) return;
        const roles = await roleRes.json();

        const containers = document.querySelectorAll(".roleDropdownContent");
        const doorIDs = [...containers].map(c => c.dataset.doorId);
        const doors = await Promise.all(doorIDs.map(id => getDoorByID(id)));
        

        containers.forEach((container, i) => {
            const currentRoles = doors[i]?.roles || [];
            container.innerHTML = "";

            roles.forEach(role => {
                const label = document.createElement("label");
                label.className = "roleItem";

                const checkbox = document.createElement("input");
                checkbox.type = "checkbox";
                checkbox.value = role.id;
                currentRoles.forEach(r => {
                     if(r == role.role_name) {
                        checkbox.checked = true;
                     }
                })

                checkbox.addEventListener("change", async () => {
                    if (checkbox.checked) {
                        await addPermission(doorIDs[i], role.id);
                    } else {
                        await removePermission(doorIDs[i], role.id);
                    }
                });

                const span = document.createElement("span");
                span.textContent = role.role_name;

                label.appendChild(checkbox);
                label.appendChild(span);
                container.appendChild(label);
            });
        });
    } catch (err) {
        console.error("Error loading roles for doors:", err);
    }
}

async function addPermission(doorID, roleID) {
    const token = localStorage.getItem("token");

    await fetch(`https://smart-lock.patohru.qzz.io/api/door/${doorID}`, {
        method: "PATCH",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ role_id: roleID })
    });
}

async function removePermission(doorID, roleID) {
    const token = localStorage.getItem("token");

    await fetch(`https://smart-lock.patohru.qzz.io/api/door/${doorID}`, {
        method: "PUT",
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json"
        },
        body: JSON.stringify({ role_id: roleID })
    });
}

function toggleRoleDropdown(button) {
    const container = button.nextElementSibling;
    const isOpen = container.classList.contains("show");

    document.querySelectorAll(".roleDropdownContent.show").forEach(el => {
        el.classList.remove("show");
    });

    if (!isOpen) container.classList.add("show");
}

document.addEventListener("click", (e) => {
    if (!e.target.closest(".roleDropdown")) {
        document.querySelectorAll(".roleDropdownContent.show").forEach(el => {
            el.classList.remove("show");
        });
    }
});

function timeFormat(isoString) {
    const date = new Date(isoString);

    const options = {
        timeZone: "Asia/Ho_Chi_Minh",
        day: "2-digit",
        month: "2-digit",
        year: "numeric",
        hour12: false
    };

    const parts = new Intl.DateTimeFormat("en-GB", options).formatToParts(date);
    const get = type => parts.find(p => p.type === type).value;

    return `${get("day")}/${get("month")}/${get("year")}`;
}

document.addEventListener("DOMContentLoaded", async () => {
    await listDoors();
    await loadRolesForDoors();
});
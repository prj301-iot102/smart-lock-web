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
            }
        );
        if(!response.ok) {
            message.textContent = "Error loading doors";
            return;
        }
        const data = await response.json();
        data.forEach(door => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${door.id}</td>
                <td>${door.door_name}</td>
                <td>${door.device_id}</td>
                <td>${door.created_at}</td>
                <td>${door.updated_at}</td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error listing doors: ", error);
    }
}

async function getDoorByID(id) {
    const token = localStorage.getItem("token");
    const doorID = id || document.getElementById("doorIdInput").value.trim();
    const message = document.getElementById("searchMsg");

    message.textContent = "";

    if(!doorID) {
        message.textContent = "Please enter Door ID";
        return;
    }

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/door/${doorID}`, 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "Door not found";
            return;
        }
        const data = await response.json();
        renderDoorInfo(data);
    } catch(error) {
        console.log("Error getting door by ID: ", error);
    }
}

async function addDoorPermission() {
    const token = localStorage.getItem("token");
    const doorID = document.getElementById("doorID").textContent.trim();
    const roleID = document.getElementById("roleIdInput").value.trim();
    const message = document.getElementById("permissionMsg");

    message.textContent = "";
    
    if(!doorID) {
        message.textContent = "No door selected";
        return;
    }
    if(!roleID) {
        message.textContent = "Please enter role ID";
        return;
    }

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/door/${doorID}`, 
            {
                method: "PATCH",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    role_id: roleID
                })
            }
        );
        if(!response.ok) {
            message.textContent = "Failed to add door permission";
            return;
        }
        message.textContent = "Add door permission successfully";
        document.getElementById("roleIdInput").value = "";
        getDoorByID(doorID);
    } catch(error) {
        console.log("Error adding door permission: ", error);
    } 
}

async function deleteDoorPermission() {
    const token = localStorage.getItem("token");
    const doorID = document.getElementById("doorID").textContent.trim();
    const roleID = document.getElementById("roleIdInput").value.trim();
    const message = document.getElementById("permissionMsg");

    message.textContent = "";
    
    if(!doorID) {
        message.textContent = "No door selected";
        return;
    }
    if(!roleID) {
        message.textContent = "Please enter role ID";
        return;
    }

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/door/${doorID}`, 
            {
                method: "PUT",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    role_id: roleID
                })
            }
        );
        if(!response.ok) {
            message.textContent = "Failed to delete door permission";
            return;
        }
        message.textContent = "Delete door permission successfully";
        document.getElementById("roleIdInput").value = "";
        getDoorByID(doorID);
    } catch(error) {
        console.log("Error deleting door permission: ", error);
    } 
}

function renderDoorInfo(door) {
    const rolesContainer = document.getElementById("roles");
    document.getElementById("doorID").textContent = door.id;
    document.getElementById("doorName").textContent = door.door_name;
    document.getElementById("deviceID").textContent = door.device_id;
    rolesContainer.innerHTML = "";
    const accessRoles = door.roles || [];
    accessRoles.forEach(role => {
        const spanTag = document.createElement("span");
        spanTag.className = "role";
        spanTag.textContent = role;
        rolesContainer.appendChild(spanTag);
    });
    document.getElementById("createdAt").textContent = door.created_at;
    document.getElementById("updatedAt").textContent = door.updated_at;
}

document.addEventListener("DOMContentLoaded", () => {
    listDoors();
});
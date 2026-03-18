async function listRoles() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("rolesTableBody");
    const message = document.getElementById("rolesMsg");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/role/",
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "Failed to get roles";
            return;
        }
        const data = await response.json();      
        data.forEach(role => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${role.id}</td>
                <td>${role.role_name}</td>
                <td>${role.created_at}</td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error listing roles:", error);
    }
}

async function searchRoleByName() {
    const token = localStorage.getItem("token");
    const name = document.getElementById("roleNameInput").value.trim();
    const table = document.getElementById("rolesTableBody");
    const message = document.getElementById("rolesMsg");

    table.innerHTML = "";
    message.textContent = "";

    if(!name) {
        message.textContent = "Please enter role name";
        return;
    }

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/role/", 
            {
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    role_name: name
                })
            }
        );
        if(!response.ok) {
            message.textContent = "Role not found";
            return;
        }
        const data = await response.json();
        data.forEach(role => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${role.id}</td>
                <td>${role.role_name}</td>
                <td>${role.created_at}</td>
            `;
            table.appendChild(tableRow);
        }); 
        document.getElementById("roleNameInput").value = "";
    } catch(error) {
        console.log("Error searching role: ", error);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    listRoles();
});
document.addEventListener("DOMContentLoaded", () => {
    listDevices();
});

async function listDevices() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("devicesTableBody");
    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/devices/", 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            throw new Error("Fail to fetch devices");
        }
        const data = await response.json();
        table.innerHTML = "";
        data.forEach(device => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${device.id}</td>
                <td>${device.device_name}</td>
                <td>${device.mac_address}</td>
                <td>
                    <select onchange="enableCreate('${device.id}', this)">
                        <option value="true" ${device.can_create ? "selected" : ""}>Enabled</option>
                        <option value="false" ${!device.can_create ? "selected" : ""}>Disabled</option>
                    </select>
                </td>
                <td>${device.created_at}</td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error loading devices: ", error);
    }
}

async function enableCreate(deviceID, currentFlag) {
    const token = localStorage.getItem("token");
    const newVal = currentFlag.value === "true";
    const message = newVal ? "Do you want to enable create of this device?" : 
                             "Do you want to disable create of this device?";
    const confirmAction = confirm(message);
    if(!confirmAction) {
        listDevices();
        return;
    }
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/devices/${deviceID}/enable`, 
            {
                method: "PATCH",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    can_create: newVal
                })
            }
        );
        if(!response.ok) {
            throw new Error("Update failed");
        }
        listDevices();
    } catch(error) {
        alert("Update failed");
        console.log(error);
        listDevices();
    }
}

async function getDeviceByID() {
    const token = localStorage.getItem("token");
    const deviceID = document.getElementById("deviceIdInput").value.trim();
    const table = document.getElementById("devicesTableBody");
    if(!deviceID) {
        listDevices();
        return;
    }
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/devices/${deviceID}`, 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            throw new Error("Device not found");
        }
        const device = await response.json();
        table.innerHTML = "";
        const tableRow = document.createElement("tr");
        tableRow.innerHTML = `
            <td>${device.id}</td>
            <td>${device.device_name}</td>
            <td>${device.mac_address}</td>
            <td>
                <select onchange="enableCreate('${device.id}', this)">
                    <option value="true" ${device.can_create ? "selected" : ""}>Enabled</option>
                    <option value="false" ${!device.can_create ? "selected" : ""}>Disabled</option>
                </select>
            </td>
            <td>${device.created_at}</td>
        `;
        table.appendChild(tableRow);
    } catch(error) {
        console.log("Error getting device by ID: ", error);
    }
}


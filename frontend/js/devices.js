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
                <td class="createbtn">
                <button 
                        class="enablebtn"
                        onclick="enableCreate('${device.id}', '${device.can_create}')" 
                        data-enabled="${device.can_create}">
                        <i class="bi bi-dot" style="font-size:20px;"></i>
                        ${device.can_create ? "Enable" : "Disable"}
                </button>
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
    const newVal = currentFlag === "true";
    const message = newVal
    ? "Unable to disable create of this device!!!" 
    : "Do you want to enable create of this device?\nAre you sure? You will not be able to disable it.";
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

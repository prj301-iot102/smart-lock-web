let employeeList = [];

document.addEventListener("DOMContentLoaded", async() => {
    await fetchEmployees();
    listNfcTags();
});

async function fetchEmployees() {
    const token = localStorage.getItem("token");
 
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/?page=1&limit=100`,
            {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                    Accept: "application/json, application/xml",
                },
            }
        );
 
        if (!response.ok) {
            console.warn("Failed to load employees");
            return;
        }
 
        const employees = await response.json();
 
        if (!employees.data || employees.data.length === 0) {
            return;
        }
 
        employeeList = employees.data;
        console.log("Employees loaded:", employeeList);
    } catch (error) {
        console.log("Error getting employees: ", error);
    }
}

function buildEmployeeOptions(selectedId) {
    let option = "";
    employeeList.forEach((e) => {
        const selected = String(e.id) === String(selectedId) ? "selected" : "";
        option += `<option value="${e.id}" ${selected}>${e.full_name}: ${e.id}</option>`;
    });
    return option;
}

async function activeNFC(id) {
    const token = localStorage.getItem("token");
    if (!confirm("Are you sure to active this NFC tag?")) {
        return;
    }
    try {
        const response = await fetch(
            `https://smart-lock.patohru.qzz.io/api/nfc/${id}`,
            {
                method: "PATCH",
                headers: {
                    Accept: "application/json, application/xml",
                    Authorization: `Bearer ${token}`,
                },
            },
        );

        const data = await response.json();

        if (response.ok) {
            alert("NFC active successfully");
            listNfcTags();
        } else {
            alert(data.message || " NFC failed");
        }
    } catch (error) {
        console.log("Error: ", error);
    }
}

async function revokeNFC(id) {
    const token = localStorage.getItem("token");
    if (!confirm("Are you sure to revoke this NFC tag?")) {
        return;
    }
    console.log(id);
    try {
        const response = await fetch(
            `https://smart-lock.patohru.qzz.io/api/nfc/${id}/revoke`,
            {
                method: "PATCH",
                headers: {
                    Accept: "application/json, application/xml",
                    Authorization: `Bearer ${token}`,
                },
            },
        );

        const data = await response.json();

        if (response.ok) {
            alert("NFC revoked successfully");
            listNfcTags();
        } else {
            alert(data.message || "Rovoke NFC failed");
        }
    } catch (error) {
        console.log("Error: ", error);
    }
}

async function listNfcTags() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("nfcTableBody");
    const nfcMessage = document.getElementById("error-msg");

    table.innerHTML = "";
    nfcMessage.textContent = "";

    try {
        const response = await fetch(
            "https://smart-lock.patohru.qzz.io/api/nfc?name=",
            {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                    Accept: "application/json, application/xml",
                },
            },
        );
        if (!response.ok) {
            nfcMessage.textContent = "Failed to load NFC tags";
            return;
        }
        const data = await response.json();
        data.forEach(nfc => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${nfc.id}</td>
                <td>${nfc.uid}</td>
                <td>
                    <select id="employeeSelect-${nfc.id}">
                        ${buildEmployeeOptions(nfc.employee_id)}
                    </select>
                </td>
                <td>${nfc.full_name}</td>
                <td>${nfc.role_name}</td>
                <td>${nfc.is_active}</td>
                <td>${timeFormat(nfc.created_at)}</td>
                <td>${timeFormat(nfc.updated_at)}</td>
                <td>
                    <div class="nfc-button">
                        <button class="revoke" onclick="revokeNFC('${nfc.id}');">
                            <i class="bi bi-trash3"></i>
                            Revoke
                        </button>
                        <button class="check" onclick="activeNFC('${nfc.id}');">
                            <i class="bi bi-check2"></i>
                            Active
                        </button>
                        <button class="check" onclick="updateEmployeeID('${nfc.id}');">
                            <i class="bi bi-floppy"></i>
                            Save
                        </button>
                    </div>
                </td>
            `;
            table.appendChild(tableRow);
        });
    } catch (error) {
        console.log("Error listing NFC Tags: ", error);
    }
}

var input = document.getElementById("nfcID");
var btn = document.getElementsByClassName("searchbtn");
var listbtn = document.getElementById("listbtn");
input.addEventListener("keyup", function (e) {
    if (e.keyCode === 13) {
        e.preventDefault();
        getNFC();
    }
});

btn.addEventListener("click", function (e) {
    if (e.target.classList.contains("searchbtn")) {
        e.preventDefault();
        getNFC();
    }
});

listbtn.addEventListener("click", function (e) {
        e.preventDefault();
        listNfcTags();
});

document.addEventListener("DOMContentLoaded", () => {
    fetchEmployees();
    listNfcTags();
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

    const parts = new Intl.DateTimeFormat("en-GB", options)
        .formatToParts(date);

    const get = type => parts.find(p => p.type === type).value;

    return `${get("day")}/${get("month")}/${get("year")}`;
}

async function updateEmployeeID(nfcId) {
    const select = document.getElementById(`employeeSelect-${nfcId}`);

    const employeeID = select.value
    console.log(employeeID)
    const token = localStorage.getItem("token");

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/nfc/${nfcId}`, 
            {
                method: "PUT", 
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    employee_id: employeeID
                })
            }
        );
        console.log(response);
        if(response.ok) {
            alert("Updated ID successfully");
        } else {
            alert("Failed to Updated ID");
            throw new Error("Failed to Updated ID");
        }
    } catch(error) {
        console.log("Failed to Updated ID: ", error);
    }; 
}


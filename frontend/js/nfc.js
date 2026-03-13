async function getNFC() {
    const id = document.getElementById("nfcID").value;
    const token = localStorage.getItem("token");
    try{
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/nfc/${id}`, {
            method: "GET",
            headers: {
                "Accept": "application/json, application/xml",
                "Authorization": `Bearer ${token}`
            }
        });
        const data = await response.json();

        console.log(data);

        if(response.ok){
            const table = document.getElementById("nfcTableBody");
            table.innerHTML = `
                <tr>
                    <td>${data.id}</td>
                    <td>${data.uid}</td>
                    <td>${data.employee_id}</td>
                    <td>${data.full_name}</td>
                    <td>${data.role_name}</td>
                    <td>${data.is_active}</td>
                    <td>${data.created_at}</td>
                    <td>${data.updated_at}</td>
                    <td>
                       <button onclick="revokeNFC('${data.id}')">Revoke</button>
                    </td>
                </tr>
            `;
        } else {
            document.getElementById("error-msg").innerText = data.message || "NFC not found";
        }
    } catch(error) {
        console.log("Cannot connect to server")
    }
}

async function revokeNFC(id) {
    const token = localStorage.getItem("token");
    if(!confirm("Are you sure to revoke this NFC tag?")) {
        return;
    }
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/nfc/${id}/revoke`, {
            method: "PATCH",
            headers: {
                "Accept": "application/json, application/xml",
                "Authorization": `Bearer ${token}`
            }
        }
        );

        const data = await response.json();

        if(response.ok) {
            alert("NFC revoked successfully");
            getNFC();
        } else {
            alert(data.message || "Rovoke NFC failed");
        }
    } catch(error) {
        console.log("Error: ", error)
    }
}

async function listNfcTags() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("nfcTableBody");
    const nfcMessage = document.getElementById("error-msg");

    table.innerHTML = "";
    nfcMessage.textContent = "";

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/nfc/",
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        }
        );
        if(!response.ok) {
            nfcMessage.textContent = "Failed to load NFC tags";
            return;
        }
        const nfcTags = await response.json();
        if(!nfcTags.data || nfcTags.data.length === 0) {
            nfcMessage.textContent = "No NFC tags found";
            return;
        }

        nfcTags.data.forEach(nfc => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${nfc.id}</td>
                <td>${nfc.uid}</td>
                <td>${nfc.employee_id}</td>
                <td>${nfc.full_name}</td>
                <td>${nfc.role_name}</td>
                <td>${nfc.is_active}</td>
                <td>${nfc.created_at}</td>
                <td>${nfc.updated_at}</td>
                <td>
                    <button onclick="revokeNFC('${nfc.id}')">Revoke</button>
                </td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error listing NFC Tags: ", error);
    }
}

document.addEventListener("DOMContentLoaded", () => {
    listNfcTags();
});

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
                    <td>${data.username}</td>
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

let currentNfcID = null;

async function getNfcByID() {
    const token = localStorage.getItem("token");
    const id = document.getElementById("nfcID").value.trim();
    if(!id) {
        alert("Please enter NFC ID");
        return;
    }
    try{
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/nfc/${id}`, 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            throw new Error("NFC not found");
        }
        const nfc = await response.json();
        renderNFCInfo(nfc);
        
    } catch(error) {
        console.log("Error getting NFC: ", error);
    }
}

function renderNFCInfo(nfc) {

    currentNfcID = nfc.id;

    document.getElementById("nfcID").textContent = nfc.id;
    document.getElementById("nfcUID").textContent = nfc.uid;
    document.getElementById("employeeID").textContent = nfc.employee_id;
    document.getElementById("fullName").textContent = nfc.full_name;
    document.getElementById("username").textContent = nfc.username;
    document.getElementById("status").textContent = nfc.is_active ? "ACTIVE" : "REVOKED";
    document.getElementById("createdAt").textContent = nfc.created_at;
    document.getElementById("updatedAt").textContent = nfc.updated_at;
}

async function revokeNFC() {
    const token = localStorage.getItem("token");
    if(!currentNfcID) {
        alert("No NFC selected");
        return;
    }
    if(!confirm("Are you sure to revoke this NFC tag?")) {
        return;
    }
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/nfc/${currentNfcID}/revoke`, 
            {
                method: "PATCH",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                }
            }
        );
        if(!response.ok) {
            throw new Error("Revoke NFC failed");
        }
        alert("NFC revoked");
        getNfcByID();    
    } catch(error) {
        console.log("Error: ", error)
    }
}


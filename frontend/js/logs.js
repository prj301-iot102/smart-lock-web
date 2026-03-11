document.addEventListener("DOMContentLoaded", () => {
    listLogs();
});

async function listLogs() {
    const table = document.getElementById("logsTableBody");

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/accessLog/", 
            {
                method: "GET", 
                headers: {
                    "Accept": "application/json, application/xml"
                }
            }
        );
        const data = await response.json();
        table.innerHTML = "";
        data.forEach(log => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${log.id}</td>
                <td>${log.created_at}</td>
                <td>${log.full_name}</td>
                <td>${log.uid}</td>
                <td>${log.status}</td>
            `;
            table.appendChild(tableRow);
        });
    } catch(error) {
        console.log("Error loading logs: ", error);
    }
}
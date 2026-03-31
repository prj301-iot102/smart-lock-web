document.addEventListener("DOMContentLoaded", () => {
    listLogs();
});

async function listLogs() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("logsTableBody");
    const totaltable = document.getElementById("totalAccess");
    const totalgranted = document.getElementById("granted");
    const totaldenied = document.getElementById("denied");
    let total = 0;
    let granted = 0;
    let denied = 0;
    try {
        const response = await fetch(
            "https://smart-lock.patohru.qzz.io/api/accessLog/",
            {
                method: "GET",
                headers: {
                    Authorization: `Bearer ${token}`,
                    Accept: "application/json, application/xml",
                },
            },
        );
        const data = await response.json();
        table.innerHTML = "";
        totaldenied.innerHTML = "";
        totalgranted.innerHTML = "";
        totaltable.innerHTML = "";

        data.forEach((log) => {
            let date = new Date(log.created_at);
            let formatDate = new Intl.DateTimeFormat("vi-VN", {
                timeZone: "Asia/Ho_Chi_Minh",
                dateStyle: "full",
                timeStyle: "short",
            }).format(date);
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${formatDate}</td>
                <td>${log.full_name}</td>
                <td>${log.uid}</td>
                <td>${log.status}</td>
            `;
            total = total + 1;
            table.appendChild(tableRow);
            if (log.status === "granted") {
                tableRow.classList.add("granted");
                granted++;
            } else {
                tableRow.classList.add("denied");
                denied++;
            }
        });
        totaldenied.innerHTML = denied;
        totalgranted.innerHTML = granted;
        totaltable.innerHTML = total;
    } catch (error) {
        console.log("Error loading logs: ", error);
    }
}

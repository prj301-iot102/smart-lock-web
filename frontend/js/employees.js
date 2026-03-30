
// let currentPage = 1;
// let currentLimit = 10;
// let totalPages = 1;

document.addEventListener("DOMContentLoaded", () => {
    listEmployees();
});

async function listEmployees() {
    const token = localStorage.getItem("token");
    const table = document.getElementById("employeesTableBody");
    const message = document.getElementById("employeeMessage");
    const pageInfo = document.getElementById("pageInfo");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/?page=${1}&limit=${10}`,
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "Failed to load employees";
            return;
        }
        const employees = await response.json();
        if(!employees.data || employees.data.length === 0) {
            message.textContent = "No employees found";
            pageInfo.textContent = "";
            return;
        }
        employees.data.forEach(employee => {
            const tableRow = document.createElement("tr");
            tableRow.innerHTML = `
                <td>${employee.id}</td>
                <td>${employee.full_name}</td>
                <td>${employee.birth}</td>
                <td>${employee.department}</td>
                <td>${employee.created_at}</td>
                <td>${employee.updated_at}</td>
            `;  
            table.appendChild(tableRow);        
        });
        currentPage = employees.page;
        currentLimit = employees.limit;
        console.log(employees.total_pages);
        totalPages = employees.total_pages;
        pageInfo.textContent = `${currentPage} / ${totalPages}`;
    } catch(error) {
        console.log("Error getting employees: ", error);
    } 
}

// let currentPage = 1;
// let currentLimit = 10;
// let totalPages = 1;

document.addEventListener("DOMContentLoaded", () => {
    listEmployees(currentPage, currentLimit);

    document.getElementById("fromDate").hidden = true;
    document.getElementById("toDate").hidden = true;
    document.getElementById("filterValue").hidden = true;
});

async function listEmployees(page = currentPage, limit = currentLimit) {
    const token = localStorage.getItem("token");
    const table = document.getElementById("employeesTableBody");
    const message = document.getElementById("employeeMessage");
    const pageInfo = document.getElementById("pageInfo");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/?page=${page}&limit=${limit}`,
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
                <td>${employee.role_name}</td>
                <td>${employee.created_at}</td>
                <td>${employee.updated_at}</td>
                <td>
                    <button onclick="openUpdateModal(
                        ${employee.id},
                        '${employee.full_name}',
                        '${employee.birth}',
                        '${employee.department}',
                        ${employee.role_id}
                    )">Update</button>
                    <button onclick="deleteEmployee(${employee.id})">Delete</button>
                </td>
            `;  
            table.appendChild(tableRow);        
        });
        currentPage = employees.page;
        currentLimit = employees.limit;
        totalPages = employees.total_pages;
        pageInfo.textContent = `${currentPage} / ${totalPages}`;
    } catch(error) {
        console.log("Error getting employees: ", error);
    } 
}

function changeLimit() {
    const selectedLimit = document.getElementById("selectedLimit");
    currentLimit = parseInt(selectedLimit.value);
    currentPage = 1;
    listEmployees(currentPage, currentLimit);
}

function nextPage() {
    if (currentPage < totalPages) {
        currentPage++;
        listEmployees(currentPage, currentLimit);
    }
}

function prevPage() {
    if (currentPage > 1) {
        currentPage--;
        listEmployees(currentPage, currentLimit);
    }
}

function searchEmployee() {
    const searchType = document.getElementById("searchType").value;
    const keyword = document.getElementById("searchInput").value.trim();

    if(!keyword) {
        alert("Please enter search value");
        return;
    }

    if(searchType === "id") {
        getEmployeeByID(keyword);
    } else {
        searchEmployeeByName(keyword, 1, currentLimit);
    }
}

async function getEmployeeByID(id) {
    const token = localStorage.getItem("token");
    const table = document.getElementById("employeesTableBody");
    const message = document.getElementById("employeeMessage");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/${id}`,
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "Employee not found";
            return;
        }
        const employee = await response.json();
        const tableRow = document.createElement("tr");
        tableRow.innerHTML = `
            <td>${employee.id}</td>
            <td>${employee.full_name}</td>
            <td>${employee.birth}</td>
            <td>${employee.department}</td>
            <td>${employee.role_name}</td>
            <td>${employee.created_at}</td>
            <td>${employee.updated_at}</td>
            <td>
                <button onclick="openUpdateModal(
                    ${employee.id},
                    '${employee.full_name}',
                    '${employee.birth}',
                    '${employee.department}',
                    ${employee.role_id}
                )">Update</button>
                <button onclick="deleteEmployee(${employee.id})">Delete</button>
            </td>
        `;
        table.appendChild(tableRow);
    } catch(error) {
        console.log("Error getting employee by ID: ", error);
    }   
}

async function searchEmployeeByName(name, page = 1, limit) {
    const token = localStorage.getItem("token");
    const table = document.getElementById("employeesTableBody");
    const message = document.getElementById("employeeMessage");

    table.innerHTML = "";
    message.textContent = "";

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/search?name=${name}&page=${page}&limit=${limit}`, 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "Failed to search employees";
        }
        const employees = await response.json();
        if(!employees.data || employees.data.length === 0){
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
                <td>${employee.role_name}</td>
                <td>${employee.created_at}</td>
                <td>${employee.updated_at}</td>
            `;
            table.appendChild(tableRow);
        });
        currentPage = employees.page;
        totalPages = employees.total_pages;
        pageInfo.textContent = `${currentPage} / ${totalPages}`;
    } catch(error) {
        console.log("Error searching employees by name: ", error);
    }
}

let updatingEmployeeId = null;

function openUpdateModal(id, name, birth, department, roleId){

    updatingEmployeeId = id;

    const updateModal = document.getElementById("updateModal");

    updateModal.style.display = "block";

    document.getElementById("newFullName").value = name;
    document.getElementById("newBirth").value = birth;
    document.getElementById("newDepartment").value = department;
    document.getElementById("newRoleId").value = roleId;
}

function closeUpdateModal(){
    document.getElementById("updateModal").style.display = "none";
    updatingEmployeeId = null;
}

async function updateEmployee() {
    const token = localStorage.getItem("token");

    const newId = document.getElementById("newId").value;
    const newFullName = document.getElementById("newFullName").value;
    const newBirth = document.getElementById("newBirth").value;
    const newDepartment = document.getElementById("newDepartment").value;
    const newRoleId = document.getElementById("newRoleId").value;
    
    if(!confirm("Are you sure to update this employee?")) {
        return;
    }

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/${updatingEmployeeId}`,
            {
                method: "PUT",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    id: newId,
                    full_name: newFullName,
                    birth: newBirth,
                    department: newDepartment,
                    role_id: newRoleId
                })
            }
        );
        const data = await response.json();
        if(!response.ok) {
            document.getElementById("updateError").textContent = data.message || data.detail || "Update failed";
        }
        alert("Employee updated successfully");
        closeUpdateModal();
        listEmployees(currentPage, currentLimit);
    } catch (error) {
        console.log("Error updating employee: ", error);
    } 
}

async function deleteEmployee(id) {
    const token = localStorage.getItem("token");
    const message = document.getElementById("employeeMessage");

    message.textContent = "";

    if(!confirm("Are you sure you want to delete this employee?")){
        return;
    }

    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/employees/${id}`,
            {
                method: "DELETE",
                headers:{
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        const data = await response.json();

        if(response.ok){

            message.textContent = "Employee deleted successfully";
            listEmployees(currentPage, currentLimit);
        }else{
            message.textContent = "Delete failed";
        }

    }catch(error){
        console.log("Delete error:", error);
    }
}

function showCreateForm() {
    document.getElementById("createEmployeeForm").classList.remove("hidden-form");
}

function hideCreateForm() {
    document.getElementById("createEmployeeForm").classList.add("hidden-form");

    document.getElementById("fullName").value = "";
    document.getElementById("department").value = "";
}

async function createEmployee() {
    const token = localStorage.getItem("token");

    const fullName = document.getElementById("fullName").value;
    const department = document.getElementById("department").value;

    if(!fullName || !department) {
        alert("Please fill all fields");
        return;
    }

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/employees/", 
            {
                method: "POST", 
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    full_name: fullName,
                    department: department
                })
            }
        );

        if(response.ok) {
            alert("Employee created successfully");
            hideCreateForm();
        } else {
            alert("Failed to create employee");
            throw new Error("Failed to create employee");
        }
    } catch(error) {
        console.log("Error creating employee: ", error);
    }; 
}

function changeFilterInput(){
    const type = document.getElementById("filterType").value;

    const from = document.getElementById("fromDate");
    const to = document.getElementById("toDate");
    const value = document.getElementById("departmentInput");

    from.hidden = true;
    to.hidden = true;
    value.hidden = true;

    if(type === "birth" || type === "created" || type === "updated"){
        from.hidden = false;
        to.hidden = false;
    }
    else if(type === "department"){
        value.type = "text";
        value.hidden = false;
    }
}

function handleFilter(){
    const type = document.getElementById("filterType").value;

    if(type === "birth"){
        filterByBirth();
    }
    else if(type === "created"){
        filterByDateAdded();
    }
    else if(type === "updated"){
        filterByDateUpdated();
    }
    else if(type === "department"){
        filterByDepartment();
    }
}

async function filterByBirth(page = 1, limit = currentLimit){

    const token = localStorage.getItem("token");

    const from = document.getElementById("fromDate").value;
    const to = document.getElementById("toDate").value;
    
    if(!from && !to){
        alert("Please select start or end date")
        return;
    }

    const params = new URLSearchParams({
        page,
        limit
    });
    
    if(from) params.append("from", from);
    if(to) params.append("to", to);

    const url = `https://smart-lock.patohru.qzz.io/api/employees/filter/birth?${params}`;
    try {
        const response = await fetch(url,
            {
                method: "GET",
                headers:{
                    "Authorization":`Bearer ${token}`,
                    "Accept":"application/json, application/xml"
            }
        });
        const data = await response.json();
        renderEmployees(data.data);
        currentPage = data.page;
        totalPages = data.total_pages;
    } catch(error) {
        console.log("Error filtering by birth: ", error);
    }
}


async function filterByDateAdded(page = 1, limit = currentLimit) {
    const token = localStorage.getItem("token");

    const from = document.getElementById("fromDate").value;
    const to = document.getElementById("toDate").value;
    
    if(!from && !to){
        alert("Please select start or end date")
        return;
    }

    const params = new URLSearchParams({
        page,
        limit
    });
    
    if(from) params.append("from", from);
    if(to) params.append("to", to);

    const url = `https://smart-lock.patohru.qzz.io/api/employees/filter/date-added?${params}`;
    try {
        const response = await fetch(url,
            {
                method: "GET",
                headers:{
                    "Authorization":`Bearer ${token}`,
                    "Accept":"application/json, application/xml"
            }
        });
        const data = await response.json();
        renderEmployees(data.data);
        currentPage = data.page;
        totalPages = data.total_pages;
    } catch(error) {
        console.log("Error filtering by date added: ", error);
    }
}

async function filterByDateUpdated(page = 1, limit = currentLimit) {
    const token = localStorage.getItem("token");

    const from = document.getElementById("fromDate").value;
    const to = document.getElementById("toDate").value;
    
    if(!from && !to){
        alert("Please select start or end date")
        return;
    }

    const params = new URLSearchParams({
        page,
        limit
    });
    
    if(from) params.append("from", from);
    if(to) params.append("to", to);

    const url = `https://smart-lock.patohru.qzz.io/api/employees/filter/date-updated?${params}`;
    try {
        const response = await fetch(url,
            {
                method: "GET",
                headers:{
                    "Authorization":`Bearer ${token}`,
                    "Accept":"application/json, application/xml"
            }
        });
        const data = await response.json();
        renderEmployees(data.data);
        currentPage = data.page;
        totalPages = data.total_pages;
    } catch(error) {
        console.log("Error filtering by date updated: ", error);
    }
}

async function filterByDepartment(page = 1, limit = currentLimit) {
    const token = localStorage.getItem("token");
     const department = document.getElementById("departmentInput").value.trim();
    const message = document.getElementById("employeeMessage");

    message.textContent = "";

    if(!department){
        message.textContent = "Please enter department";
        return;
    }

    const params = new URLSearchParams({
        department,
        page,
        limit
    });

    const url = `https://smart-lock.patohru.qzz.io/api/employees/filter/department?${params}`;

    try{

        const response = await fetch(url,{
            method: "GET",
            headers:{
                "Authorization":`Bearer ${token}`,
                "Accept":"application/json, application/xml"
            }
        });
        if(!response.ok){
            message.textContent = "Failed to filter employees by department";
            return;
        }
        const data = await response.json();
        if(!data.data || data.data.length === 0){
            message.textContent = "No employees found";
            return;
        }
        renderEmployees(data.data);
        currentPage = data.page;
        totalPages = data.total_pages;
    }catch(error){
        console.log("Filter department error:", error);
    }
}

function renderEmployees(employees){

    const table = document.getElementById("employeesTableBody");
    const message = document.getElementById("employeeMessage");

    table.innerHTML = "";
    message.textContent = "";

    if(!employees || employees.length === 0){
        message.textContent = "No employees found";
        return;
    }

    employees.forEach(employee => {

        const tableRow = document.createElement("tr");

        tableRow.innerHTML = `
            <td>${employee.id}</td>
            <td>${employee.full_name}</td>
            <td>${employee.birth}</td>
            <td>${employee.department}</td>
            <td>${employee.role_name}</td>
            <td>${employee.created_at}</td>
            <td>${employee.updated_at}</td>
            <td>
                <button onclick="openUpdateModal(
                    ${employee.id},
                    '${employee.full_name}',
                    '${employee.birth}',
                    '${employee.department}',
                    ${employee.role_id}
                )">Update</button>
                <button onclick="deleteEmployee(${employee.id})">
                    Delete
                </button>
            </td>
        `;
        table.appendChild(tableRow);
    });
}


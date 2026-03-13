
async function getUserByID() {
    const token = localStorage.getItem("token");
    const id = document.getElementById("userIdInput").value.trim();
    const table = document.getElementById("userTableBody");
    const message = document.getElementById("userMessage");

    table.innerHTML = "";
    message.textContent = "";

    if(!id) {
        alert("Please input user ID");
        return;
    }
    try {
        const response = await fetch(`https://smart-lock.patohru.qzz.io/api/users/${id}`, 
            {
                method: "GET",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml"
                }
            }
        );
        if(!response.ok) {
            message.textContent = "User not found";
            return;
        }
        const user = await response.json();
        const tableRow = document.createElement("tr");
        tableRow.innerHTML = `
            <td>${user.id}</td>
            <td>${user.full_name}</td>
            <td>${user.username}</td>
            <td>${user.role_name}</td>
            <td>${user.created_at}</td>
        `;
        table.appendChild(tableRow);
    } catch(error) {
        console.log("Error getting user: ", error);
    } 
}

function showCreateForm() {
    document.getElementById("createUserForm").classList.remove("hidden-form");
}

function hideCreateForm() {
    document.getElementById("createUserForm").classList.add("hidden-form");

    document.getElementById("fullName").value = "";
    document.getElementById("department").value = "";
    document.getElementById("username").value = "";
    document.getElementById("password").value = "";
}

async function createUser() {
    const token = localStorage.getItem("token");

    const fullName = document.getElementById("fullName").value;
    const department = document.getElementById("department").value;
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    if(!fullName || !department || !username || !password) {
        alert("Please fill all fields");
        return;
    }

    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/users/create", 
            {
                method: "POST", 
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    full_name: fullName,
                    department: department,
                    username: username,
                    password: password
                })
            }
        );

        if(response.ok) {
            alert("User created successfully");
            hideCreateForm();
        } else {
            alert("Failed to create user");
            throw new Error("Failed to create user");
        }
    } catch(error) {
        console.log("Error creating user: ", error);
    }; 
}

async function updatePassword() {
    const token = localStorage.getItem("token");
    const newPassword = document.getElementById("newPassword").value.trim();
    if(!newPassword) {
        alert("Please enter new password");
        return;
    }
    const confirmUpdate = confirm("Are you sure you want to change your password?");
    if(!confirmUpdate) {
        return;
    }
    try {
        const response = await fetch("https://smart-lock.patohru.qzz.io/api/users/update", 
            {
                method: "PATCH",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Accept": "application/json, application/xml",
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    password: newPassword
                })
            }
        );
        if(response.ok) {
            alert("Password updated successfully");
            document.getElementById("newPassword").value = "";
        } else {
            alert("Failed to update password");
        }
    } catch(error) {
        console.log("Error updating password: ", error);
    }
}

var input = document.getElementById("userIdInput");
input.addEventListener("keyup", function(e) {
    if (e.keyCode = 13) {
        getUserByID;
    }
});
document.addEventListener("click", function(e) {
    if (e.target.id === "logoutBtn") {
        const isConfirm = confirm("Are you sure?");

        if (isConfirm) {
            localStorage.removeItem("token");
            window.location.replace("./index.html");
        }
    }
});
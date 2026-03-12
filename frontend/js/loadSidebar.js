document.addEventListener("DOMContentLoaded", () => {

fetch("components/sidebar.html")
.then(res => res.text())
.then(html => {

    const sidebarElement = document.getElementById("sidebar");
    sidebarElement.innerHTML = html;

    const links = sidebarElement.querySelectorAll(".sidebar-menu a");

    const currentPage = window.location.pathname.split("/").pop();

    links.forEach(link => {

        link.classList.remove("active");

        const href = link.getAttribute("href");

        if (currentPage && href && currentPage.endsWith(href)) {
            link.classList.add("active");
        }

    });

});

});
function updateUserUI(email) {
    const loginBtn = document.getElementById("login-btn");
    const userInfoDiv = document.getElementById("user-info");
    const userEmailSpan = document.getElementById("user-email");

    if (email) {
        userEmailSpan.textContent = email;
        userInfoDiv.style.display = "flex";
        loginBtn.style.display = "none";
    } else {
        userEmailSpan.textContent = "";
        userInfoDiv.style.display = "none";
        loginBtn.style.display = "block";
    }
}

function fetchUserData() {
    fetch("/user")
        .then(response => {
            if (response.ok) {
                return response.json();
            }
            throw new Error("Не удалось получить данные пользователя");
        })
        .then(data => {
            updateUserUI(data.email);
        })
        .catch(() => {
            updateUserUI(null);
        });
}

document.addEventListener("DOMContentLoaded", fetchUserData);


function loginUser() {
    const email = document.getElementById("login-email").value.trim();
    const password = document.getElementById("login-password").value.trim();

    if (!email || !password) {
        alert("Все поля обязательны для заполнения.");
        return;
    }

    fetch("/login", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
    })
        .then((res) => res.json())
        .then((data) => {
            if (data.email) {
                updateUserUI(data.email);
                closeModal('login-modal');
            } else {
                alert("Ошибка входа");
            }
        })
        .catch((err) => {
            console.error("Ошибка при входе:", err);
        });
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}
document.addEventListener("DOMContentLoaded", function () {
    const compareList = new Set();

    document.querySelector(".university-list").addEventListener("click", function (e) {
        const target = e.target;

        // Обработка кнопок "Сравнить"
        if (target.classList.contains("compare-btn")) {
            const universityID = target.closest(".university-card").dataset.universityId;

            if (compareList.has(universityID)) {
                compareList.delete(universityID);
                target.textContent = "Сравнить";
            } else {
                compareList.add(universityID);
                target.textContent = "К сравнению";
            }

            console.log("Выбрано для сравнения:", Array.from(compareList));

            const compareLink = document.querySelector(".navigation-menu a[href='/compare']");
            compareLink.addEventListener("click", function (e) {
                e.preventDefault();

                if (compareList.size === 0) {
                    alert("Выберите хотя бы один университет для сравнения.");
                    return;
                }

                const queryString = Array.from(compareList).join(",");
                window.location.href = `/compare?ids=${queryString}`;
            });
        }

        // Обработка кнопок "В Избранное"
        if (target.classList.contains("detail-btn")) {
            const universityID = parseInt(target.dataset.universityId);

            if (isNaN(universityID)) return;

            const isFavorite = target.classList.contains("in-favorites");
            const url = isFavorite ? "/remove/favorites" : "/add/favorites";

            fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ university_id: universityID })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.status === "success") {
                        if (isFavorite) {
                            target.classList.remove("in-favorites");
                            target.textContent = "В Избранное";
                        } else {
                            target.classList.add("in-favorites");
                            target.textContent = "В Избранном";
                        }
                    }
                })
                .catch(err => {
                    console.error("Ошибка при обновлении избранного:", err.message);
                    alert(`Ошибка: ${err.message}`);
                });
        }
    });
});

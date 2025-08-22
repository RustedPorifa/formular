class VariantsManager {
  constructor() {
    this.variantsHolder = document.getElementById("variants-holder");
    this.currentClass = this.extractClassFromUrl();
    this.apiUrl = this.buildApiUrl();
    this.init();
  }

  // Извлекаем класс из URL (например: http://127.0.0.1:5050/api/show/7 -> "7")
  extractClassFromUrl() {
    const pathSegments = window.location.pathname.split("/");

    // Ищем сегмент с классом после /api/show/
    const showIndex = pathSegments.indexOf("show");
    if (showIndex !== -1 && showIndex + 1 < pathSegments.length) {
      return pathSegments[showIndex + 1];
    }

    // Альтернативный вариант: последний сегмент URL
    return pathSegments[pathSegments.length - 1];
  }

  // Строим URL для API запроса
  buildApiUrl() {
    if (this.currentClass && this.currentClass.match(/^\d+$/)) {
      return `/api/variants/${this.currentClass}`;
    }
    return "/api/variants"; // Все варианты если класс не указан
  }

  init() {
    this.updatePageTitle();
    this.checkAuthAndUpdateProfile();
    this.loadVariants();
    this.setupEventListeners();
  }

  // Проверка авторизации и обновление профиля
  async checkAuthAndUpdateProfile() {
    try {
      const response = await fetch("/verify", {
        method: "GET",
        credentials: "include",
      });

      if (response.ok) {
        const data = await response.json();
        if (data.verify === "true") {
          this.updateProfileButton();
        }
      }
    } catch (error) {
      console.error("Ошибка при проверке аутентификации:", error);
    }
  }

  // Обновление кнопки профиля
  updateProfileButton() {
    const authBtn = document.getElementById("auth-profile-btn");
    if (!authBtn) return;

    authBtn.onclick = function () {
      window.location = "/user/profile";
    };
    authBtn.textContent = "Профиль";
  }

  // Обновляем заголовок страницы с учетом класса
  updatePageTitle() {
    if (this.currentClass && this.currentClass.match(/^\d+$/)) {
      document.title = `Formular - ${this.currentClass} класс`;

      // Также можно обновить заголовок на странице
      const pageHeader = document.querySelector("h1");
      if (pageHeader) {
        pageHeader.textContent = `Варианты для ${this.currentClass} класса`;
      }
    }
  }

  async loadVariants() {
    try {
      this.showLoading();
      console.log(
        `Загрузка вариантов для класса: ${this.currentClass || "все"}`,
      );

      const response = await fetch(this.apiUrl);

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      // Обрабатываем оба формата ответа: с success полем и без
      let variants = data.variants || data;

      if (Array.isArray(variants)) {
        this.renderVariants(variants);
      } else {
        this.showError("Неверный формат данных от сервера");
      }
    } catch (error) {
      console.error("Ошибка загрузки вариантов:", error);
      this.showError("Не удалось загрузить варианты. Проверьте соединение.");
    }
  }

  renderVariants(variants) {
    if (!variants || variants.length === 0) {
      this.showEmptyState();
      return;
    }

    // Фильтруем варианты по классу (на случай если API вернул все)
    const filteredVariants = this.currentClass
      ? variants.filter((variant) => variant.class === this.currentClass)
      : variants;

    if (filteredVariants.length === 0) {
      this.showEmptyStateForClass();
      return;
    }

    const variantsHTML = filteredVariants
      .map((variant) => this.createVariantCard(variant))
      .join("");

    this.variantsHolder.innerHTML = `
            <div class="variants-container">
                ${variantsHTML}
            </div>
        `;

    this.setupCardInteractions();
  }

  createVariantCard(variant) {
    const solvedIcon = variant.solved ? "✅" : "⏳";
    const hasVideo = variant.videoFilePath && variant.videoFilePath !== "";

    return `
            <div class="variant-card" data-uuid="${variant.uuid}">
                <div class="variant-header">
                    <h3 class="variant-title">${this.escapeHtml(variant.name)}</h3>
                    <span class="variant-badge">${variant.class} класс</span>
                </div>

                <div class="variant-meta">
                    <span class="meta-item">
                        <span class="meta-icon">📚</span>
                        ${this.escapeHtml(variant.subject)}
                    </span>
                    <span class="meta-item">
                        <span class="meta-icon">${solvedIcon}</span>
                        ${variant.solved ? "Решен" : "В процессе"}
                    </span>
                </div>

                ${
                  variant.description
                    ? `
                    <p class="variant-description">
                        ${this.escapeHtml(variant.description)}
                    </p>
                `
                    : ""
                }

                <div class="variant-actions">
                    <a href="/api/get-variant/${variant.uuid}"
                       class="action-btn action-btn-primary"
                       target="_blank">
                        👁️ Просмотреть вариант
                    </a>
                </div>
            </div>
        `;
  }

  showEmptyStateForClass() {
    this.variantsHolder.innerHTML = `
            <div class="empty-state">
                <h3>📝 Варианты не найдены</h3>
                <p>Для ${this.currentClass} класса пока нет доступных вариантов</p>
                <button onclick="window.location='/variants'" class="action-btn action-btn-primary" style="margin-top: 1rem;">
                    👀 Посмотреть все варианты
                </button>
            </div>
        `;
  }

  setupCardInteractions() {
    const cards = this.variantsHolder.querySelectorAll(".variant-card");

    cards.forEach((card) => {
      card.addEventListener("click", (e) => {
        if (!e.target.closest(".action-btn")) {
          const uuid = card.dataset.uuid;
          this.showVariantDetails(uuid);
        }
      });
    });
  }

  showVariantDetails(uuid) {
    // Открываем вариант в новой вкладке
    window.open(`/api/get-variant/${uuid}`, "_blank");
  }

  showLoading() {
    this.variantsHolder.innerHTML = `
            <div class="loading">
                <div class="loading-spinner"></div>
                <p style="margin-top: 1rem;">Загрузка вариантов для ${this.currentClass || "всех"} классов...</p>
            </div>
        `;
  }

  showError(message) {
    this.variantsHolder.innerHTML = `
            <div class="error-message">
                <h3>😕 Ошибка загрузки</h3>
                <p>${this.escapeHtml(message)}</p>
                <button onclick="location.reload()" class="action-btn action-btn-primary" style="margin-top: 1rem;">
                    🔄 Попробовать снова
                </button>
            </div>
        `;
  }

  showEmptyState() {
    this.variantsHolder.innerHTML = `
            <div class="empty-state">
                <h3>📝 Варианты не найдены</h3>
                <p>Пока нет доступных вариантов для отображения</p>
            </div>
        `;
  }

  escapeHtml(text) {
    const div = document.createElement("div");
    div.textContent = text;
    return div.innerHTML;
  }

  setupEventListeners() {
    document.addEventListener("keydown", (e) => {
      if (e.key === "r" && e.ctrlKey) {
        e.preventDefault();
        this.loadVariants();
      }
    });
  }
}

// Инициализация при загрузке страницы
document.addEventListener("DOMContentLoaded", () => {
  new VariantsManager();
});

// Глобальные функции
window.VariantsManager = {
  refresh: function () {
    new VariantsManager().loadVariants();
  },

  filterByClass: function (className) {
    window.location.href = `/api/variants/${className}`;
  },
};

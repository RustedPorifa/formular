// Инициализация даты
document.addEventListener("DOMContentLoaded", () => {
  const now = new Date();
  const options = {
    day: "2-digit",
    month: "long",
    year: "numeric",
  };
  document.getElementById("currentDate").textContent = now.toLocaleDateString(
    "ru-RU",
    options,
  );
  document.getElementById("currentYear").textContent = now.getFullYear();
  // Инициализация темы
  const savedTheme = localStorage.getItem("theme");
  if (savedTheme === "dark") {
    document.body.classList.add("dark-theme");
  }
  // Инициализация первого блока
  initVariantBlock(document.querySelector(".variant-block"));
});

// Переключение между разделами
const navItems = document.querySelectorAll(".nav-item");
const emailSection = document.getElementById("emailSection");
const uploadSection = document.getElementById("uploadSection");
const contentTitle = document.getElementById("contentTitle");
navItems.forEach((item) => {
  item.addEventListener("click", function () {
    const section = this.getAttribute("data-section");
    // Удаляем активный класс у всех элементов
    navItems.forEach((i) => i.classList.remove("active"));
    // Добавляем активный класс текущему элементу
    this.classList.add("active");
    // Переключаем разделы
    if (section === "email") {
      emailSection.classList.add("active");
      uploadSection.style.display = "none";
      contentTitle.textContent = "Email рассылка";
    } else if (section === "upload") {
      emailSection.classList.remove("active");
      uploadSection.style.display = "block";
      contentTitle.textContent = "Массовая загрузка вариантов";
    }
  });
});
// По умолчанию активен раздел загрузки
document.querySelector('[data-section="upload"]').classList.add("active");

// Зависимость видов материалов от класса
const subjectsMap = {
  5: ["ВПР", "Математика"],
  6: ["ВПР", "Математика"],
  7: ["ВПР", "Алгебра", "Геометрия"],
  8: ["ВПР", "Алгебра", "Геометрия"],
  9: ["ОГЭ", "Алгебра", "Геометрия"],
  10: ["ВПР", "Алгебра", "Геометрия", "Стереометрия"],
  11: ["ЕГЭ", "Алгебра", "Геометрия", "Стереометрия"],
};
function initClassSelect(select) {
  select.addEventListener("change", function () {
    const variantBlock = this.closest(".variant-block");
    const subjectSelect = variantBlock.querySelector(".subject-select");
    const selectedClass = this.value;
    if (selectedClass) {
      subjectSelect.disabled = false;
      subjectSelect.innerHTML = "";
      // Добавляем опции в зависимости от класса
      subjectsMap[selectedClass].forEach((subject) => {
        const option = document.createElement("option");
        option.value = subject.toLowerCase();
        option.textContent = subject;
        subjectSelect.appendChild(option);
      });
    } else {
      subjectSelect.disabled = true;
      subjectSelect.innerHTML =
        '<option value="">Сначала выберите класс</option>';
    }
  });
}

// Инициализация блока варианта
function initVariantBlock(block) {
  // Инициализация выбора класса
  const classSelect = block.querySelector(".class-select");
  initClassSelect(classSelect);

  // Загрузка PDF файла
  const pdfInput = block.querySelector(".pdf-input");
  const pdfUpload = block.querySelector(".pdf-upload");
  const pdfFileName = block.querySelector(".pdf-upload .file-name");
  const pdfProgress = block.querySelector(".pdf-upload .progress-bar");
  const pdfProgressFill = block.querySelector(".pdf-upload .progress-fill");

  pdfInput.addEventListener("change", function (e) {
    if (this.files.length > 0) {
      const file = this.files[0];
      // Проверка расширения файла
      if (!file.name.toLowerCase().endsWith(".pdf")) {
        alert("Пожалуйста, выберите файл в формате PDF");
        this.value = "";
        pdfFileName.textContent = "Файл не выбран";
        pdfUpload.style.borderColor = "";
        return;
      }
      // Проверка размера файла (максимум 20MB)
      if (file.size > 20 * 1024 * 1024) {
        alert("Файл слишком большой. Максимальный размер - 20MB");
        this.value = "";
        pdfFileName.textContent = "Файл не выбран";
        pdfUpload.style.borderColor = "";
        return;
      }
      pdfFileName.textContent = file.name;
      pdfUpload.style.borderColor = "var(--success)";
    }
  });

  // Загрузка видео файла
  const videoInput = block.querySelector(".video-input");
  const videoUpload = block.querySelector(".video-upload");
  const videoFileName = block.querySelector(".video-upload .file-name");
  const videoProgress = block.querySelector(".video-upload .progress-bar");
  const videoProgressFill = block.querySelector(".video-upload .progress-fill");

  videoInput.addEventListener("change", function (e) {
    if (this.files.length > 0) {
      const file = this.files[0];
      // Проверка расширения файла
      if (!file.name.toLowerCase().endsWith(".mp4")) {
        alert("Пожалуйста, выберите файл в формате MP4");
        this.value = "";
        videoFileName.textContent = "Файл не выбран";
        videoUpload.style.borderColor = "";
        return;
      }
      // Проверка размера файла (максимум 10гб)
      if (file.size > 10000 * 1024 * 1024) {
        alert("Файл слишком большой. Максимальный размер - 100MB");
        this.value = "";
        videoFileName.textContent = "Файл не выбран";
        videoUpload.style.borderColor = "";
        return;
      }
      videoFileName.textContent = file.name;
      videoUpload.style.borderColor = "var(--success)";
    }

    // Удаление блока
    const removeBtn = block.querySelector(".remove-block");
    removeBtn.addEventListener("click", function () {
      block.remove();
      updateBlockNumbers();
    });
  });
}

// Добавление нового блока варианта
const addBlockBtn = document.getElementById("addBlockBtn");
const variantBlocks = document.getElementById("variantBlocks");
addBlockBtn.addEventListener("click", function () {
  const blockCount = variantBlocks.children.length + 1;
  const newBlock = document.createElement("div");
  newBlock.className = "variant-block";
  newBlock.innerHTML = `
    <button class="remove-block">×</button>
    <div class="block-header">
        <div class="block-number">${blockCount}</div>
        <h3 class="block-title">Вариант #${blockCount}</h3>
    </div>
    <form class="variant-form">
        <div class="form-grid">
            <div class="form-group">
                <label>Название варианта *</label>
                <input type="text" class="variant-name" placeholder="Введите название" required>
            </div>
            <div class="form-group">
                <label>Описание (опционально)</label>
                <input type="text" class="variant-desc" placeholder="Краткое описание">
            </div>
        </div>
        <div class="form-grid">
            <div class="form-group">
                <label>Класс *</label>
                <select class="class-select" required>
                    <option value="">Выберите класс</option>
                    <option value="5">5 класс</option>
                    <option value="6">6 класс</option>
                    <option value="7">7 класс</option>
                    <option value="8">8 класс</option>
                    <option value="9">9 класс</option>
                    <option value="10">10 класс</option>
                    <option value="11">11 класс</option>
                </select>
            </div>
            <div class="form-group">
                <label>Вид материала *</label>
                <select class="subject-select" disabled required>
                    <option value="">Сначала выберите класс</option>
                </select>
            </div>
        </div>
        <div class="form-grid">
            <div class="form-group">
                <label>PDF файл *</label>
                <div class="file-upload pdf-upload">
                    <i class="fas fa-file-pdf"></i>
                    <p>Перетащите PDF файл сюда или кликните для выбора</p>
                    <div class="file-name">Файл не выбран</div>
                    <input type="file" class="file-input pdf-input" accept=".pdf">
                    <div class="progress-bar">
                        <div class="progress-fill"></div>
                    </div>
                </div>
            </div>
            <div class="form-group">
                <label>Видео (MP4, опционально)</label>
                <div class="file-upload video-upload">
                    <i class="fas fa-video"></i>
                    <p>Перетащите видеофайл сюда или кликните для выбора</p>
                    <div class="file-name">Файл не выбран</div>
                    <input type="file" class="file-input video-input" accept=".mp4">
                    <div class="progress-bar">
                        <div class="progress-fill"></div>
                    </div>
                </div>
            </div>
        </div>
        <div class="solved-checkbox">
            <input type="checkbox" class="solved-checkbox-input">
            <label>Решенный вариант</label>
        </div>
    </form>
`;
  variantBlocks.appendChild(newBlock);
  initVariantBlock(newBlock);
  updateBlockNumbers();
});

// Обновление номеров блоков
function updateBlockNumbers() {
  const blocks = document.querySelectorAll(".variant-block");
  blocks.forEach((block, index) => {
    const number = index + 1;
    block.querySelector(".block-number").textContent = number;
    block.querySelector(".block-title").textContent = `Вариант #${number}`;
    // Показываем кнопку удаления для всех блоков, кроме первого
    if (number > 1) {
      block.querySelector(".remove-block").style.display = "flex";
    } else {
      block.querySelector(".remove-block").style.display = "none";
    }
  });
}

// Отправка всех вариантов
const uploadAllBtn = document.getElementById("uploadAllBtn");
const uploadStatus = document.getElementById("uploadStatus");
uploadAllBtn.addEventListener("click", function () {
  const blocks = document.querySelectorAll(".variant-block");
  let valid = true;
  let errorMessages = [];

  // Сбросим предыдущие статусы
  uploadStatus.innerHTML = "";
  blocks.forEach((block) => {
    block.classList.remove("error");
  });

  // Проверяем каждый блок
  blocks.forEach((block, index) => {
    const nameInput = block.querySelector(".variant-name");
    const classSelect = block.querySelector(".class-select");
    const subjectSelect = block.querySelector(".subject-select");
    const pdfInput = block.querySelector(".pdf-input");
    const pdfFileName = block.querySelector(".pdf-upload .file-name");
    const blockNumber = index + 1;
    let blockErrors = [];

    // Проверка обязательных полей
    if (!nameInput.value) {
      blockErrors.push("отсутствует название");
    }
    if (!classSelect.value) {
      blockErrors.push("не выбран класс");
    }
    if (!subjectSelect.value) {
      blockErrors.push("не выбран вид материала");
    }
    if (pdfFileName.textContent === "Файл не выбран") {
      blockErrors.push("не выбран PDF файл");
    }

    // Проверка PDF файла
    if (pdfInput.files.length > 0) {
      const file = pdfInput.files[0];
      if (!file.name.toLowerCase().endsWith(".pdf")) {
        blockErrors.push("неверный формат PDF файла");
      }
    }

    // Проверка видео файла
    const videoInput = block.querySelector(".video-input");
    if (videoInput.files.length > 0) {
      const file = videoInput.files[0];
      if (!file.name.toLowerCase().endsWith(".mp4")) {
        blockErrors.push("неверный формат видео файла");
      }
    }

    if (blockErrors.length > 0) {
      valid = false;
      block.classList.add("error");
      const variantName = nameInput.value || `Вариант #${blockNumber}`;
      errorMessages.push(
        `<strong>${variantName}</strong>: ${blockErrors.join(", ")}`,
      );
    }
  });

  if (!valid) {
    let errorHtml = `<div class="status-error">Обнаружены ошибки:<br>`;
    errorMessages.forEach((msg) => {
      errorHtml += `• ${msg}<br>`;
    });
    errorHtml += `</div>`;
    uploadStatus.innerHTML = errorHtml;
    return;
  }

  // Создаем FormData для отправки
  const formData = new FormData();

  // Собираем данные из всех блоков
  blocks.forEach((block, index) => {
    const name = block.querySelector(".variant-name").value;
    const desc = block.querySelector(".variant-desc").value;
    const classVal = block.querySelector(".class-select").value;
    const subjectVal = block.querySelector(".subject-select").value;
    const solved = block.querySelector(".solved-checkbox-input").checked;
    const pdfFile = block.querySelector(".pdf-input").files[0];
    const videoFile = block.querySelector(".video-input").files[0];

    // Добавляем данные варианта в formData
    formData.append(`variants[${index}][name]`, name);
    formData.append(`variants[${index}][description]`, desc);
    formData.append(`variants[${index}][class]`, classVal);
    formData.append(`variants[${index}][subject]`, subjectVal);
    formData.append(`variants[${index}][solved]`, solved);

    // Добавляем файлы
    if (pdfFile) {
      formData.append(`variants[${index}][pdf]`, pdfFile);
    }
    if (videoFile) {
      formData.append(`variants[${index}][video]`, videoFile);
    }
  });

  // Блокируем кнопку во время загрузки
  uploadAllBtn.disabled = true;
  uploadAllBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Загрузка...';
  uploadStatus.innerHTML = `<div class="status-info">Начата загрузка ${blocks.length} вариантов...</div>`;

  // Отправляем данные на сервер
  fetch("/admin/upload-variants", {
    method: "POST",
    body: formData,
  })
    .then((response) => {
      if (!response.ok) {
        return response.json().then((err) => {
          throw new Error(err.message || "Ошибка сети");
        });
      }
      return response.json();
    })
    .then((data) => {
      if (data.success) {
        uploadStatus.innerHTML = `<div class="status-success">Успешно загружено ${data.uploaded} вариантов из ${blocks.length}!</div>`;

        // Очистка формы (оставляем первый блок, удаляем остальные)
        blocks.forEach((block, index) => {
          if (index === 0) {
            // Оставляем первый блок
            block.querySelector("form").reset();
            block.querySelector(".pdf-upload .file-name").textContent =
              "Файл не выбран";
            block.querySelector(".video-upload .file-name").textContent =
              "Файл не выбран";
            block.querySelector(".pdf-upload").style.borderColor = "";
            block.querySelector(".video-upload").style.borderColor = "";
            block.querySelector(".pdf-upload .progress-bar").style.display =
              "none";
            block.querySelector(".video-upload .progress-bar").style.display =
              "none";
          } else {
            // Удаляем остальные блоки
            block.remove();
          }
        });
        updateBlockNumbers();
      } else {
        throw new Error(data.message || "Ошибка при загрузке");
      }
    })
    .catch((error) => {
      console.error("Ошибка:", error);
      uploadStatus.innerHTML = `<div class="status-error">Ошибка: ${error.message}</div>`;
    })
    .finally(() => {
      uploadAllBtn.disabled = false;
      uploadAllBtn.innerHTML =
        '<i class="fas fa-cloud-upload-alt"></i> Загрузить все варианты';
    });
});

// Обработка email рассылки
const emailForm = document.getElementById("emailForm");
emailForm.addEventListener("submit", function (e) {
  e.preventDefault();
  const subject = document.getElementById("emailSubject").value;
  const content = document.getElementById("emailContent").value;
  if (!subject || !content) {
    alert("Пожалуйста, заполните тему и содержание письма!");
    return;
  }

  // Отправляем данные на сервер
  fetch("/api/send-email", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      subject: subject,
      content: content,
    }),
  })
    .then((response) => response.json())
    .then((data) => {
      if (data.success) {
        alert(`Рассылка "${subject}" успешно отправлена!`);
        emailForm.reset();
      } else {
        alert(`Ошибка при отправке: ${data.message}`);
      }
    })
    .catch((error) => {
      console.error("Ошибка:", error);
      alert("Произошла ошибка при отправке рассылки");
    });
});

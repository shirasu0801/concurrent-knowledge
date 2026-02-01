// Concurrent Knowledge - Admin Panel JavaScript
// Handles question list, filtering, and deletion

let currentPage = 1;
let currentGenreFilter = '';
let currentDifficultyFilter = '';

document.addEventListener('DOMContentLoaded', function() {
    // Load initial question list
    loadQuestions();

    // Filter button
    document.getElementById('filterButton').addEventListener('click', function() {
        currentGenreFilter = document.getElementById('genreFilter').value;
        currentDifficultyFilter = document.getElementById('difficultyFilter').value;
        currentPage = 1;
        loadQuestions();
    });
});

async function loadQuestions(page = 1) {
    const questionList = document.getElementById('questionList');
    questionList.innerHTML = '<p>読み込み中...</p>';

    try {
        // Build query parameters
        const params = new URLSearchParams({
            page: page,
            pageSize: 20
        });

        if (currentGenreFilter) {
            params.append('genreId', currentGenreFilter);
        }

        if (currentDifficultyFilter) {
            params.append('difficulty', currentDifficultyFilter);
        }

        // Fetch questions
        const response = await fetch('/api/questions?' + params.toString());
        const data = await response.json();

        if (!data.success) {
            questionList.innerHTML = '<p class="error">問題の読み込みに失敗しました</p>';
            return;
        }

        const questions = data.data.questions;
        const total = data.data.total;
        const totalPages = data.data.totalPages;

        // Display questions
        if (questions.length === 0) {
            questionList.innerHTML = '<p class="no-items">問題が見つかりませんでした</p>';
            return;
        }

        questionList.innerHTML = questions.map(q => `
            <div class="question-item">
                <div>
                    <div style="margin-bottom: 0.5rem;">
                        <span class="badge">${getDifficultyBadge(q.difficulty)}</span>
                        <span class="badge">${q.genreId}</span>
                    </div>
                    <div style="font-weight: 600; margin-bottom: 0.25rem;">${q.questionText}</div>
                    <div style="font-size: 0.9rem; color: var(--text-secondary);">
                        正解: ${q.correctAnswer} | 基本点: ${q.basePoints}
                    </div>
                </div>
                <div style="display: flex; gap: 0.5rem;">
                    <a href="/admin/questions/${q.id}/edit" class="btn btn-small btn-secondary">編集</a>
                    <button class="btn btn-small btn-danger delete-button" data-id="${q.id}" data-text="${escapeHtml(q.questionText)}">削除</button>
                </div>
            </div>
        `).join('');

        // Add delete event listeners
        document.querySelectorAll('.delete-button').forEach(button => {
            button.addEventListener('click', function() {
                const id = this.dataset.id;
                const text = this.dataset.text;
                deleteQuestion(id, text);
            });
        });

        // Display pagination
        displayPagination(page, totalPages);

    } catch (error) {
        console.error('Error loading questions:', error);
        questionList.innerHTML = '<p class="error">問題の読み込みに失敗しました</p>';
    }
}

function getDifficultyBadge(difficulty) {
    const colors = {
        '初級': '#3fb950',
        '中級': '#d29922',
        '上級': '#f85149'
    };
    const color = colors[difficulty] || '#6e7681';
    return `<span style="background-color: ${color}; color: white; padding: 0.25rem 0.5rem; border-radius: 4px; font-size: 0.85rem;">${difficulty}</span>`;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function displayPagination(currentPage, totalPages) {
    const pagination = document.getElementById('pagination');

    if (totalPages <= 1) {
        pagination.innerHTML = '';
        return;
    }

    let html = '';

    // Previous button
    if (currentPage > 1) {
        html += `<button class="btn btn-small btn-secondary" onclick="loadQuestions(${currentPage - 1})">前へ</button>`;
    }

    // Page numbers
    for (let i = 1; i <= totalPages; i++) {
        if (i === currentPage) {
            html += `<button class="btn btn-small btn-primary">${i}</button>`;
        } else if (i === 1 || i === totalPages || Math.abs(i - currentPage) <= 2) {
            html += `<button class="btn btn-small btn-secondary" onclick="loadQuestions(${i})">${i}</button>`;
        } else if (Math.abs(i - currentPage) === 3) {
            html += `<span>...</span>`;
        }
    }

    // Next button
    if (currentPage < totalPages) {
        html += `<button class="btn btn-small btn-secondary" onclick="loadQuestions(${currentPage + 1})">次へ</button>`;
    }

    pagination.innerHTML = html;
    window.currentPage = currentPage;
}

async function deleteQuestion(id, questionText) {
    if (!confirm(`この問題を削除してもよろしいですか？\n\n"${questionText}"`)) {
        return;
    }

    try {
        const response = await fetch(`/api/questions/${id}`, {
            method: 'DELETE'
        });

        const data = await response.json();

        if (data.success) {
            alert('問題を削除しました');
            loadQuestions(window.currentPage || 1);
        } else {
            alert('削除に失敗しました: ' + data.error);
        }
    } catch (error) {
        console.error('Error deleting question:', error);
        alert('削除に失敗しました');
    }
}

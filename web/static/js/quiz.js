// Concurrent Knowledge - Quiz Controller
// Handles quiz gameplay, answer submission, and result tracking

class QuizController {
    constructor() {
        this.quizData = null;
        this.currentQuestionIndex = 0;
        this.answers = [];
        this.timer = null;
        this.startTime = Date.now();
    }

    // Initialize quiz from sessionStorage
    async init() {
        // Get quiz data from sessionStorage
        const quizDataStr = sessionStorage.getItem('currentQuiz');
        if (!quizDataStr) {
            alert('クイズデータが見つかりません');
            window.location.href = '/';
            return;
        }

        this.quizData = JSON.parse(quizDataStr);

        // Get current question index and answers
        this.currentQuestionIndex = parseInt(sessionStorage.getItem('currentQuestionIndex') || '0');
        const answersStr = sessionStorage.getItem('quizAnswers');
        this.answers = answersStr ? JSON.parse(answersStr) : [];

        // Display current question
        this.displayQuestion(this.currentQuestionIndex);
    }

    // Display a question
    displayQuestion(index) {
        if (index >= this.quizData.totalQuestions) {
            this.showResults();
            return;
        }

        const question = this.quizData.questions[index];

        // Update question number
        document.getElementById('questionNumber').textContent = index + 1;
        document.getElementById('totalQuestions').textContent = this.quizData.totalQuestions;

        // Update question text
        document.getElementById('questionText').textContent = question.questionText;

        // Create option buttons
        const optionsContainer = document.getElementById('optionsContainer');
        optionsContainer.innerHTML = '';

        Object.keys(question.options).forEach(key => {
            const button = document.createElement('button');
            button.className = 'option-button';
            button.textContent = `${key}. ${question.options[key]}`;
            button.dataset.answer = key;

            button.addEventListener('click', () => this.selectAnswer(key));

            optionsContainer.appendChild(button);
        });

        // Start timer
        this.startTimer(this.quizData.timePerQuestion);

        // Hide answer feedback
        document.getElementById('answerFeedback').style.display = 'none';
    }

    // Start countdown timer
    startTimer(seconds) {
        let remainingTime = seconds;
        const timerElement = document.getElementById('timer');

        // Clear any existing timer
        if (this.timer) {
            clearInterval(this.timer);
        }

        // Update display immediately
        timerElement.textContent = remainingTime;
        timerElement.className = 'quiz-timer';

        // Start countdown
        this.timer = setInterval(() => {
            remainingTime--;
            timerElement.textContent = remainingTime;

            // Update timer color based on remaining time
            if (remainingTime <= 10) {
                timerElement.className = 'quiz-timer danger';
            } else if (remainingTime <= 20) {
                timerElement.className = 'quiz-timer warning';
            }

            // Time's up
            if (remainingTime <= 0) {
                clearInterval(this.timer);
                this.submitAnswer('', 0); // Submit with no answer
            }
        }, 1000);

        // Store remaining time for later use
        this.remainingTime = remainingTime;
    }

    // Select an answer
    async selectAnswer(answer) {
        // Stop timer
        if (this.timer) {
            clearInterval(this.timer);
        }

        // Get remaining time from timer display
        const timerElement = document.getElementById('timer');
        const remainingTime = parseInt(timerElement.textContent);

        // Highlight selected option
        document.querySelectorAll('.option-button').forEach(btn => {
            btn.classList.remove('selected');
            if (btn.dataset.answer === answer) {
                btn.classList.add('selected');
            }
            btn.disabled = true;
        });

        // Submit answer
        await this.submitAnswer(answer, remainingTime);
    }

    // Submit answer to server
    async submitAnswer(answer, remainingTime) {
        const question = this.quizData.questions[this.currentQuestionIndex];

        try {
            const response = await fetch('/api/quiz/submit', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    quizId: this.quizData.quizId,
                    questionId: question.id,
                    answer: answer,
                    remainingTime: remainingTime
                })
            });

            const data = await response.json();

            if (data.success) {
                const result = data.data;

                // Store answer result with question details
                const answerRecord = {
                    questionId: question.id,
                    questionText: question.questionText,
                    yourAnswer: answer,
                    correct: result.correct,
                    correctAnswer: result.correctAnswer,
                    explanation: result.explanation,
                    pointsEarned: result.pointsEarned
                };

                this.answers.push(answerRecord);
                sessionStorage.setItem('quizAnswers', JSON.stringify(this.answers));

                // Save to review list if incorrect
                if (!result.correct) {
                    ReviewStorage.add(answerRecord);
                }

                // Show feedback
                this.showFeedback(result);
            } else {
                alert('回答の送信に失敗しました: ' + data.error);
            }
        } catch (error) {
            console.error('Error submitting answer:', error);
            alert('回答の送信に失敗しました');
        }
    }

    // Show answer feedback
    showFeedback(result) {
        const feedbackElement = document.getElementById('answerFeedback');
        const resultElement = document.getElementById('feedbackResult');
        const explanationElement = document.getElementById('feedbackExplanation');
        const scoreElement = document.getElementById('feedbackScore');

        // Show feedback
        feedbackElement.style.display = 'block';

        // Set result text and class
        if (result.correct) {
            resultElement.textContent = '正解！';
            resultElement.className = 'feedback-result correct';
        } else {
            resultElement.textContent = '不正解';
            resultElement.className = 'feedback-result incorrect';
        }

        // Set explanation
        explanationElement.textContent = `正解: ${result.correctAnswer} - ${result.explanation}`;

        // Set score
        if (result.correct) {
            scoreElement.textContent = `獲得ポイント: ${result.pointsEarned}点 (基本: ${result.calculation.basePoints}, 倍率: ${result.calculation.difficultyMultiplier}x, 時間ボーナス: ${result.calculation.timeBonus})`;
        } else {
            scoreElement.textContent = '獲得ポイント: 0点';
        }

        // Setup next button
        const nextButton = document.getElementById('nextButton');
        nextButton.onclick = () => this.nextQuestion();
    }

    // Move to next question
    nextQuestion() {
        this.currentQuestionIndex++;
        sessionStorage.setItem('currentQuestionIndex', this.currentQuestionIndex.toString());

        if (this.currentQuestionIndex < this.quizData.totalQuestions) {
            this.displayQuestion(this.currentQuestionIndex);
        } else {
            this.showResults();
        }
    }

    // Show results page
    showResults() {
        window.location.href = '/result';
    }
}

// Initialize quiz on page load
document.addEventListener('DOMContentLoaded', function() {
    const controller = new QuizController();
    controller.init();
});

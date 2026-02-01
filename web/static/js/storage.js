// Concurrent Knowledge - LocalStorage Management
// Manages review list (incorrect answers) in browser localStorage

class ReviewStorage {
    static KEY = 'concurrentKnowledge_reviewList';

    // Get all review items
    static getAll() {
        try {
            const data = localStorage.getItem(this.KEY);
            return data ? JSON.parse(data) : [];
        } catch (error) {
            console.error('Error reading from localStorage:', error);
            return [];
        }
    }

    // Add an incorrect answer to review list
    static add(incorrectAnswer) {
        try {
            const current = this.getAll();

            // Check if question already exists (avoid duplicates)
            const exists = current.some(item => item.questionId === incorrectAnswer.questionId);
            if (exists) {
                // Update existing entry
                const index = current.findIndex(item => item.questionId === incorrectAnswer.questionId);
                current[index] = {
                    ...incorrectAnswer,
                    timestamp: Date.now()
                };
            } else {
                // Add new entry
                current.push({
                    ...incorrectAnswer,
                    timestamp: Date.now()
                });
            }

            localStorage.setItem(this.KEY, JSON.stringify(current));
        } catch (error) {
            console.error('Error writing to localStorage:', error);
            if (error.name === 'QuotaExceededError') {
                alert('復習リストの容量が上限に達しました。古いアイテムを削除してください。');
            }
        }
    }

    // Remove a specific question from review list
    static remove(questionId) {
        try {
            const current = this.getAll();
            const filtered = current.filter(item => item.questionId !== questionId);
            localStorage.setItem(this.KEY, JSON.stringify(filtered));
        } catch (error) {
            console.error('Error removing from localStorage:', error);
        }
    }

    // Clear all review items
    static clear() {
        try {
            localStorage.removeItem(this.KEY);
        } catch (error) {
            console.error('Error clearing localStorage:', error);
        }
    }

    // Cleanup old entries (older than 30 days)
    static cleanup() {
        try {
            const current = this.getAll();
            const thirtyDaysAgo = Date.now() - (30 * 24 * 60 * 60 * 1000);

            const filtered = current.filter(item => {
                return item.timestamp && item.timestamp > thirtyDaysAgo;
            });

            if (filtered.length !== current.length) {
                localStorage.setItem(this.KEY, JSON.stringify(filtered));
                console.log(`Cleaned up ${current.length - filtered.length} old review items`);
            }
        } catch (error) {
            console.error('Error during cleanup:', error);
        }
    }

    // Get count of review items
    static count() {
        return this.getAll().length;
    }
}

// Auto cleanup on page load
document.addEventListener('DOMContentLoaded', function() {
    ReviewStorage.cleanup();
});

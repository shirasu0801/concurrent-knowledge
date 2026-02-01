// Concurrent Knowledge - Timer Utility
// Simple timer utility for quiz gameplay

class Timer {
    constructor(duration, onTick, onComplete) {
        this.duration = duration;
        this.remaining = duration;
        this.onTick = onTick;
        this.onComplete = onComplete;
        this.interval = null;
    }

    start() {
        this.remaining = this.duration;

        // Call onTick immediately
        if (this.onTick) {
            this.onTick(this.remaining);
        }

        // Start interval
        this.interval = setInterval(() => {
            this.remaining--;

            if (this.onTick) {
                this.onTick(this.remaining);
            }

            if (this.remaining <= 0) {
                this.stop();
                if (this.onComplete) {
                    this.onComplete();
                }
            }
        }, 1000);
    }

    stop() {
        if (this.interval) {
            clearInterval(this.interval);
            this.interval = null;
        }
    }

    getRemaining() {
        return this.remaining;
    }
}

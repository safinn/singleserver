const sloganBadge = document.querySelector('[data-slogan-rotator]');
if (sloganBadge) {
    const sloganText = sloganBadge.querySelector('.tagline-text');
    const slogans = JSON.parse(sloganBadge.dataset.slogans);
    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    let sloganIndex = 0;
    let timer = null;

    const advance = () => {
        sloganIndex = (sloganIndex + 1) % slogans.length;

        if (reduceMotion) {
            sloganText.textContent = slogans[sloganIndex];
            return;
        }

        sloganBadge.classList.add('is-changing');
        window.setTimeout(() => {
            sloganText.textContent = slogans[sloganIndex];
        }, 320);
        window.setTimeout(() => {
            sloganBadge.classList.remove('is-changing');
        }, 700);
    };

    const schedule = () => {
        if (timer) {
            window.clearInterval(timer);
        }
        timer = window.setInterval(advance, 3000);
    };

    sloganBadge.addEventListener('click', () => {
        if (sloganBadge.classList.contains('is-changing')) {
            return;
        }
        advance();
        schedule();
    });

    schedule();
}

document.querySelectorAll('.copy-btn').forEach((btn) => {
    btn.addEventListener('click', async () => {
        const codeBlock = btn.previousElementSibling;
        const text = codeBlock.innerText.trim();
        const originalIcon = btn.innerHTML;

        try {
            await navigator.clipboard.writeText(text);
            btn.innerHTML = '<svg viewBox="0 0 24 24" fill="currentColor" style="color: #00ff88;"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/></svg>';
            setTimeout(() => {
                btn.innerHTML = originalIcon;
            }, 1600);
        } catch (err) {
            console.error('Failed to copy:', err);
        }
    });
});

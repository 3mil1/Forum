import React, {useEffect, useRef} from 'react';
import styles from './matrix404.module.css'

function MatrixBackground404({timeout = 50}) {
    const canvas = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        const context = canvas?.current?.getContext('2d');

        const width = document.body.offsetWidth;
        const height = document.body.offsetHeight;
        if (canvas?.current?.width) {
            canvas.current.width = width;
        }
        if (canvas?.current?.height) {
            canvas.current.height = height;
        }

        if (context) {
            context.fillStyle = '#000';
        }
        context?.fillRect(0, 0, width, height);

        const columns = Math.floor(width / 20) + 1;
        const yPositions = Array.from({length: columns}).fill(0);

        context!.fillStyle = '#000';
        context!.fillRect(0, 0, width, height);

        const matrixEffect = () => {
            context!.fillStyle = '#0001';
            context!.fillRect(0, 0, width, height);

            context!.fillStyle = '#0f0';
            context!.font = '15pt monospace';

            yPositions.forEach((y: any, index) => {
                const text = String.fromCharCode(Math.random() * 128);
                const x = index * 20;
                context!.fillText(text, x, y);

                if (y > 100 + Math.random() * 10000) {
                    yPositions[index] = 0;
                } else {
                    yPositions[index] = y + 20;
                }
            });
        };

        const interval = setInterval(matrixEffect, timeout);
        return () => {
            clearInterval(interval);
        };
    }, [canvas, timeout]);

    return (
        <>
            <div className={styles.t404}>404</div>
            <div
                style={{
                    background: '#000000',
                    overflow: 'hidden',
                    position: 'fixed',
                    height: '100%',
                    width: '100%',
                    zIndex: -1,
                    left: '0',
                    top: '0',
                }}
            >
                <canvas
                    ref={canvas}
                />
            </div>
        </>

    );
}

export default MatrixBackground404;
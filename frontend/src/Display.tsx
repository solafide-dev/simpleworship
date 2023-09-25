import { Context, useEffect, useRef } from "react"

type SlideStyle = {
    background: {
        color: string
        image?: string
        margin?: number
    },
    text: {
        color: string
        size: number
        horizontalAlignment?: string
        verticalAlignment?: string
        gap?: number
        uppercase?: boolean
        margin?: number
        font?: string
        fontWeight?: number
        background?: {
            color?: string | CanvasGradient | CanvasPattern
            fullWidth?: boolean
            margin?: number
        }
    }
}

function Display({ lyrics }: { lyrics: string }) {
    const canvasRef = useRef<HTMLCanvasElement>(null)
    let lastSlide: ImageBitmap

    const style = {
        background: {
            color: "#be123c",
            // video: 'ParticleSpinBlueHD.mp4',
            // image: 'ParticleSpinPinkSlowHD.jpg',
            margin: 4
        },
        transition: null, // TODO: Quick patch for build. @michaelpanik fix this.
        text: {
            color: "#111111",
            size: 3,
            horizontalAlignment: "center",
            verticalAlignment: "center",
            gap: 5,
            uppercase: true,
            margin: 5,
            font: 'Inter Variable',
            fontWeight: 600,
            background: {
                color: "#ffffff",
                fullWidth: true,
                margin: 1
            }
        }
    }
    // const style: SlideStyle = {
    //     background: {
    //         color: "#be123c",
    //         image: 'ParticleSpinPinkSlowHD.jpg',
    //         margin: 4
    //     },
    //     text: {
    //         color: "#000000",
    //         size: 4,
    //         horizontalAlignment: "center",
    //         verticalAlignment: "center",
    //         gap: 5,
    //         uppercase: true,
    //         margin: 5,
    //         font: 'Inter Variable',
    //         fontWeight: 500,
    //         background: {
    //             color: "#ffffff",
    //             fullWidth: true,
    //             margin: 1
    //         }
    //     }
    // }

    const drawText = (ctx: OffscreenCanvasRenderingContext2D) => {
        const fontSize = style.text.size * 25
        const margin = style.text.margin ? style.text.margin * 25 : 0

        ctx.font = `${style.text.fontWeight} ${fontSize}px ${style.text.font}`;
        lyrics.split('\n').map((line, i) => {
            if (style.text.uppercase) { line = line.toUpperCase() }

            let x = margin
            let y = margin + (fontSize * (i + 1))

            if (style.text.verticalAlignment == "center") {
                y = (ctx.canvas.height / 2) + (fontSize * i) + ((style.text.gap || 0) * 25 * i)
            }

            if (style.text.horizontalAlignment == "center") {
                ctx.textAlign = "center"
                x = ctx.canvas.width / 2
            }

            if (style.text?.background?.color) {
                ctx.fillStyle = style.text.background.color
                if (style.text.background.fullWidth) {
                    ctx.fillRect(0 + margin, y - fontSize - (style.text.background.margin || 0) * 25, ctx.canvas.width - (margin * 2), fontSize + (style.text.background.margin || 0) * 75)
                }
            }

            ctx.fillStyle = style.text.color
            ctx.fillText(line, x + (style.text.background?.margin || 0) * 25, y, ctx.canvas.width - margin)
        })
    }

    const createOffscreenCanvas = (width: number = 1920, height: number = 1080): OffscreenCanvas => {
        const canvas = new OffscreenCanvas(width, height)
        const ctx = canvas.getContext('2d')
        if (ctx) {
            // if (style.background.image) {
            //     const image = new Image()
            //     image.src = style.background.image
            //     image.onload = () => {
            //         ctx.drawImage(image, 0, 0, 1920, 1080)
            //         drawText(ctx)
            //     }
            // } else {
            ctx.fillStyle = style.background.color
            ctx.fillRect(0, 0, ctx.canvas.width, ctx.canvas.height)
            drawText(ctx)
            // }
        }
        return canvas
    }

    useEffect(() => {
        const ctx = canvasRef?.current?.getContext('2d')

        if (ctx) {
            if (lastSlide) { ctx.drawImage(lastSlide, 0, 0) }

            (async () => {
                const offscreenCanvas = createOffscreenCanvas(ctx.canvas.width, ctx.canvas.height)
                const imageData = await offscreenCanvas.transferToImageBitmap()
                lastSlide = imageData

                if (style.transition) {
                    let i = 0

                    const loop = () => {
                        i++

                        if (i < 100) {
                            ctx.globalAlpha = i / 100
                            ctx.drawImage(imageData, 0, 0)
                            window.requestAnimationFrame(loop)
                        }
                    }

                    window.requestAnimationFrame(loop)
                } else {
                    ctx.drawImage(imageData, 0, 0)
                }
            })()
        }

    }, [lyrics])


    return <div style={{ width: 1920, height: 1080, maxWidth: '100%' }}>
        {/* <video>
            <source src={`./${style.background.video}`} type="video/mp4" />
        </video> */}
        <canvas ref={canvasRef} width={1920} height={1080} className="max-w-full"></canvas>
    </div>
}

export default Display
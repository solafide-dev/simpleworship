if (!window.OffscreenCanvas) {
    window.OffscreenCanvas = class OffscreenCanvas {
        canvas

        constructor(width, height) {
            this.canvas = document.createElement("canvas");
            this.canvas.width = width;
            this.canvas.height = height;

            this.canvas.convertToBlob = () => {
                return new Promise(resolve => {
                    this.canvas.toBlob(resolve);
                });
            };

            this.canvas.transferToImageBitmap = async () => {
                return await createImageBitmap(this.canvas)
            }

            return this.canvas;
        }
    };
}
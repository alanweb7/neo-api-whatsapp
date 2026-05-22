declare module "qrcode" {
  interface QRCodeToDataURLOptions {
    errorCorrectionLevel?: string;
    type?: string;
    quality?: number;
    margin?: number;
    width?: number;
    color?: {
      dark?: string;
      light?: string;
    };
  }

  function toDataURL(
    text: string,
    options?: QRCodeToDataURLOptions
  ): Promise<string>;

  export { toDataURL };
}

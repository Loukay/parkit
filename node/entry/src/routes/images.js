import { Router } from 'express';
import { uploadImage } from '../services/imageUploader.js';

const router = Router();

router.post('/', uploadImage);

export { router as imagesRouter };

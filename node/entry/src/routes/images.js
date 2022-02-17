import { Router } from 'express';

const router = Router();

router.post('/', function (req, res, next) {
  res.send('respond with a resoursce');
});

export { router as imagesRouter };

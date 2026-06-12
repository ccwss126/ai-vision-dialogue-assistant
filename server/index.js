import 'dotenv/config'
import express from 'express'
import OpenAI from 'openai'

const app = express()
const port = process.env.PORT || 3001

app.use(express.json({ limit: '10mb' }))

app.post('/api/analyze', async (req, res) => {
  try {
    const { imageBase64, question } = req.body
    const apiKey = process.env.OPENAI_API_KEY

    if (!apiKey) {
      return res.status(500).json({ error: '缺少 OPENAI_API_KEY' })
    }

    if (!imageBase64) {
      return res.status(400).json({ error: '缺少 imageBase64' })
    }

    const openai = new OpenAI({ apiKey })
    const imageUrl = imageBase64.startsWith('data:')
      ? imageBase64
      : `data:image/png;base64,${imageBase64}`

    const response = await openai.responses.create({
      model: process.env.OPENAI_MODEL || 'gpt-4o-mini',
      input: [
        {
          role: 'user',
          content: [
            {
              type: 'input_text',
              text: question || '请分析当前画面。',
            },
            {
              type: 'input_image',
              image_url: imageUrl,
            },
          ],
        },
      ],
    })

    res.json({ answer: response.output_text })
  } catch (error) {
    console.error(error)
    res.status(500).json({ error: 'AI 分析失败' })
  }
})

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`)
})

---
created_at: "2026-02-28T03:15:04Z"
description: "Implement `Orchestra.Desktop/Services/VisionService.cs` — AI-powered image understanding combining Claude Vision and Windows built-in OCR.\n\n**Windows OCR (`Windows.Media.Ocr`):**\n```csharp\nvar engine = OcrEngine.TryCreateFromUserProfileLanguages();\nOcrResult result = await engine.RecognizeAsync(softwareBitmap);\nstring text = result.Text;\n```\n\n**Claude Vision:** send captured `SoftwareBitmap` as base64 PNG in `ai_prompt` tool call with `provider=claude`, image attachment\n\n**`VisionService` API (ai.vision — 6 tools):**\n- `analyze_image` — Claude Vision: describe UI, code, diagrams in screenshot\n- `extract_text` — Windows OCR first (fast, offline), Claude Vision fallback for handwriting/complex layouts  \n- `find_elements` — locate specific UI elements by description\n- `compare_images` — before/after diff description\n- `describe_screen` — full screen narration for accessibility\n- `extract_data` — pull structured data (tables, forms, charts) from screenshot\n\n**Integration with ai.screenshot:** pipeline: capture → analyze → result card in chat\n\n**`ImageAttachment`** — `ChatMessage.Attachments` list, rendered as `Image` control in chat bubble with expand-on-click\n\n**Platform:** Desktop (OCR Win10+, Vision requires internet)"
id: FEAT-SKO
priority: P1
project_id: orchestra-win
status: backlog
title: AI Vision plugin — image analysis + Windows OCR
updated_at: "2026-02-28T03:15:04Z"
version: 0
---

# AI Vision plugin — image analysis + Windows OCR

Implement `Orchestra.Desktop/Services/VisionService.cs` — AI-powered image understanding combining Claude Vision and Windows built-in OCR.

**Windows OCR (`Windows.Media.Ocr`):**
```csharp
var engine = OcrEngine.TryCreateFromUserProfileLanguages();
OcrResult result = await engine.RecognizeAsync(softwareBitmap);
string text = result.Text;
```

**Claude Vision:** send captured `SoftwareBitmap` as base64 PNG in `ai_prompt` tool call with `provider=claude`, image attachment

**`VisionService` API (ai.vision — 6 tools):**
- `analyze_image` — Claude Vision: describe UI, code, diagrams in screenshot
- `extract_text` — Windows OCR first (fast, offline), Claude Vision fallback for handwriting/complex layouts  
- `find_elements` — locate specific UI elements by description
- `compare_images` — before/after diff description
- `describe_screen` — full screen narration for accessibility
- `extract_data` — pull structured data (tables, forms, charts) from screenshot

**Integration with ai.screenshot:** pipeline: capture → analyze → result card in chat

**`ImageAttachment`** — `ChatMessage.Attachments` list, rendered as `Image` control in chat bubble with expand-on-click

**Platform:** Desktop (OCR Win10+, Vision requires internet)

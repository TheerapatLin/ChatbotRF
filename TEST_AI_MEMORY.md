# วิธีทดสอบ AI Memory (Conversation History)

## วิธีที่ 1: ทดสอบด้วย curl

### ขั้นตอนที่ 1: บอกข้อมูลส่วนตัว

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สวัสดีครับ ผมชื่อสมชาย อายุ 25 ปี ทำงานเป็น Software Engineer ที่กรุงเทพ ชอบเล่นเกม Valorant และดูอนิเมะ One Piece",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": false`
- `"history_count": 0`
- AI ทักทายและตอบรับข้อมูล

---

### ขั้นตอนที่ 2: ถามชื่อ

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "ผมชื่ออะไรครับ?",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": true`
- `"history_count": 2` (user message + AI response)
- AI ตอบ: "ชื่อสมชายครับ" หรือคล้ายๆ กัน

---

### ขั้นตอนที่ 3: ถามอายุ

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "อายุเท่าไหร่?",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": true`
- `"history_count": 4`
- AI ตอบ: "อายุ 25 ปีครับ"

---

### ขั้นตอนที่ 4: ถามงาน

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "ทำงานอะไร?",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": true`
- `"history_count": 6`
- AI ตอบ: "ทำงานเป็น Software Engineer ครับ"

---

### ขั้นตอนที่ 5: ถามงานอดิเรก

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "ชอบทำอะไรยามว่าง?",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": true`
- `"history_count": 8`
- AI ตอบ: "ชอบเล่นเกม Valorant และดูอนิเมะ One Piece ครับ"

---

### ขั้นตอนที่ 6: ถามสรุปทั้งหมด

```bash
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "สรุปข้อมูลของผมให้หน่อย",
    "session_id": "test_memory_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- `"history_used": true`
- `"history_count": 10` (ถึง limit แล้ว)
- AI สรุปข้อมูลทั้งหมด: ชื่อ, อายุ, งาน, งานอดิเรก

---

## วิธีที่ 2: ทดสอบหลายหัวข้อ

### Test Case: ซื้อของออนไลน์

```bash
# 1. บอกว่าต้องการซื้อของ
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "ผมต้องการซื้อ iPhone 15 Pro Max สีดำ 256GB",
    "session_id": "shopping_001",
    "use_history": true
  }'

# 2. เพิ่มสินค้าอื่น
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "เพิ่ม AirPods Pro 2 ด้วย",
    "session_id": "shopping_001",
    "use_history": true
  }'

# 3. ถามว่าสั่งอะไรบ้าง
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "ผมสั่งอะไรไปบ้างแล้ว?",
    "session_id": "shopping_001",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:** AI ตอบรายการสินค้าทั้งหมด (iPhone + AirPods)

---

## วิธีที่ 3: ทดสอบด้วย Postman

### Collection: AI Memory Test

**Request 1: Introduction**
```
POST http://localhost:3001/api/chat
Content-Type: application/json

{
  "message": "Hi, my name is John. I'm learning Vue.js and Go programming.",
  "session_id": "{{session_id}}",
  "use_history": true
}
```

**Request 2: Ask Name**
```
POST http://localhost:3001/api/chat
Content-Type: application/json

{
  "message": "What is my name?",
  "session_id": "{{session_id}}",
  "use_history": true
}
```

**Request 3: Ask Skills**
```
POST http://localhost:3001/api/chat
Content-Type: application/json

{
  "message": "What programming languages am I learning?",
  "session_id": "{{session_id}}",
  "use_history": true
}
```

---

## วิธีที่ 4: ทดสอบ Limit 10 Messages

### สคริปต์ทดสอบ (Bash)

```bash
#!/bin/bash

SESSION_ID="test_limit_$(date +%s)"
API_URL="http://localhost:3001/api/chat"

echo "Testing 10 message limit with session: $SESSION_ID"
echo ""

# ส่งข้อความ 15 ข้อความ
for i in {1..15}; do
  echo "=== Message $i ==="

  curl -s -X POST $API_URL \
    -H "Content-Type: application/json" \
    -d "{
      \"message\": \"This is message number $i\",
      \"session_id\": \"$SESSION_ID\",
      \"use_history\": true
    }" | jq -r '.history_count'

  echo ""
done

echo ""
echo "=== Final Test: Ask to summarize ==="

curl -s -X POST $API_URL \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Can you list all message numbers I sent?\",
    \"session_id\": \"$SESSION_ID\",
    \"use_history\": true
  }" | jq -r '.reply'
```

**ผลลัพธ์ที่คาดหวัง:**
- ข้อความที่ 11-15 ควรมี `history_count: 10` (max)
- AI ควรเห็นเฉพาะ message 6-15 (10 ข้อความล่าสุด)

---

## วิธีที่ 5: ทดสอบหลาย Session

```bash
# Session 1
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "My name is Alice",
    "session_id": "session_alice",
    "use_history": true
  }'

# Session 2 (คนละคน)
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "My name is Bob",
    "session_id": "session_bob",
    "use_history": true
  }'

# ถาม Session 1
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is my name?",
    "session_id": "session_alice",
    "use_history": true
  }'

# ถาม Session 2
curl -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "What is my name?",
    "session_id": "session_bob",
    "use_history": true
  }'
```

**ผลลัพธ์ที่คาดหวัง:**
- Session 1 ตอบ: "Alice"
- Session 2 ตอบ: "Bob"
- ไม่มีการปนกัน!

---

## วิธีที่ 6: ทดสอบแบบ Automated (Node.js)

```javascript
const axios = require('axios');

async function testAIMemory() {
  const API_URL = 'http://localhost:3001/api/chat';
  const sessionId = `test_${Date.now()}`;

  console.log(`Testing with session: ${sessionId}\n`);

  // Test cases
  const tests = [
    {
      message: "My name is Sarah and I'm 28 years old",
      expectedHistoryUsed: false,
      expectedCount: 0
    },
    {
      message: "What is my name?",
      expectedHistoryUsed: true,
      expectedCount: 2,
      expectedReply: /Sarah/i
    },
    {
      message: "How old am I?",
      expectedHistoryUsed: true,
      expectedCount: 4,
      expectedReply: /28/
    }
  ];

  for (let i = 0; i < tests.length; i++) {
    const test = tests[i];
    console.log(`\n=== Test ${i + 1} ===`);
    console.log(`Message: "${test.message}"`);

    try {
      const response = await axios.post(API_URL, {
        message: test.message,
        session_id: sessionId,
        use_history: true
      });

      const data = response.data;

      // Check history_used
      if (data.history_used === test.expectedHistoryUsed) {
        console.log(`✅ history_used: ${data.history_used}`);
      } else {
        console.log(`❌ history_used: expected ${test.expectedHistoryUsed}, got ${data.history_used}`);
      }

      // Check history_count
      if (data.history_count === test.expectedCount) {
        console.log(`✅ history_count: ${data.history_count}`);
      } else {
        console.log(`❌ history_count: expected ${test.expectedCount}, got ${data.history_count}`);
      }

      // Check reply content (if specified)
      if (test.expectedReply) {
        if (test.expectedReply.test(data.reply)) {
          console.log(`✅ reply contains expected content`);
        } else {
          console.log(`❌ reply doesn't contain expected content`);
          console.log(`   Reply: "${data.reply}"`);
        }
      }

      console.log(`Reply: "${data.reply}"`);

    } catch (error) {
      console.error(`❌ Error:`, error.message);
    }

    // Wait a bit between requests
    await new Promise(resolve => setTimeout(resolve, 1000));
  }
}

testAIMemory();
```

**รันด้วย:**
```bash
npm install axios
node test_ai_memory.js
```

---

## วิธีที่ 7: ทดสอบบน Frontend (Vue.js)

```vue
<template>
  <div class="ai-memory-test">
    <h2>AI Memory Test</h2>

    <div class="test-info">
      <p>Session ID: <code>{{ sessionId }}</code></p>
      <button @click="startNewSession">New Session</button>
    </div>

    <div class="test-cases">
      <div v-for="(test, index) in testCases" :key="index" class="test-case">
        <h3>Test {{ index + 1 }}</h3>
        <p><strong>Message:</strong> {{ test.message }}</p>
        <button @click="runTest(test, index)" :disabled="test.running">
          {{ test.running ? 'Running...' : 'Run Test' }}
        </button>

        <div v-if="test.result" class="result">
          <p>
            <span :class="test.result.historyUsed ? 'pass' : 'fail'">
              History Used: {{ test.result.historyUsed }}
            </span>
          </p>
          <p>
            <span :class="test.result.historyCount === test.expectedCount ? 'pass' : 'fail'">
              History Count: {{ test.result.historyCount }}
              (Expected: {{ test.expectedCount }})
            </span>
          </p>
          <p><strong>AI Reply:</strong> {{ test.result.reply }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      sessionId: `test_${Date.now()}`,
      testCases: [
        {
          message: 'My name is John, I love programming in Go',
          expectedHistoryUsed: false,
          expectedCount: 0,
          running: false,
          result: null
        },
        {
          message: 'What is my name?',
          expectedHistoryUsed: true,
          expectedCount: 2,
          running: false,
          result: null
        },
        {
          message: 'What programming language do I love?',
          expectedHistoryUsed: true,
          expectedCount: 4,
          running: false,
          result: null
        }
      ]
    };
  },
  methods: {
    startNewSession() {
      this.sessionId = `test_${Date.now()}`;
      this.testCases.forEach(test => {
        test.result = null;
        test.running = false;
      });
    },

    async runTest(test, index) {
      test.running = true;

      try {
        const response = await fetch('http://localhost:3001/api/chat', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            message: test.message,
            session_id: this.sessionId,
            use_history: true
          })
        });

        const data = await response.json();

        test.result = {
          historyUsed: data.history_used,
          historyCount: data.history_count,
          reply: data.reply
        };

      } catch (error) {
        test.result = {
          error: error.message
        };
      } finally {
        test.running = false;
      }
    }
  }
};
</script>

<style scoped>
.test-case {
  border: 1px solid #ddd;
  padding: 1rem;
  margin: 1rem 0;
  border-radius: 8px;
}

.pass {
  color: green;
  font-weight: bold;
}

.fail {
  color: red;
  font-weight: bold;
}

.test-info {
  background: #f5f5f5;
  padding: 1rem;
  margin-bottom: 1rem;
  border-radius: 8px;
}

code {
  background: #eee;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: monospace;
}
</style>
```

---

## สรุปวิธีทดสอบ

### ✅ ทดสอบพื้นฐาน (Basic Tests)
1. บอกข้อมูล → ถามข้อมูลเดียวกัน (ชื่อ, อายุ, งาน)
2. เช็ค `history_used: true` และ `history_count` เพิ่มขึ้น

### ✅ ทดสอบความจำหลายข้อมูล (Multiple Info)
1. บอกข้อมูลหลายอย่างในครั้งเดียว
2. ถามทีละข้อมูล
3. ถามสรุปทั้งหมด

### ✅ ทดสอบ Limit (10 Messages)
1. ส่งข้อความมากกว่า 10 ข้อความ
2. เช็คว่า `history_count` ไม่เกิน 10
3. ถามเกี่ยวกับข้อความเก่าสุด (ควรไม่จำ)

### ✅ ทดสอบหลาย Session (Isolation)
1. สร้าง 2 sessions ต่างกัน
2. บอกข้อมูลคนละอย่าง
3. ถามใน session แต่ละอัน (ต้องไม่ปน!)

### ✅ ทดสอบ Error Cases
1. ไม่ส่ง `use_history` → ไม่ควรมี history
2. ส่ง `session_id` ใหม่ → `history_count: 0`
3. `session_id` ผิด → ควรไม่มี history

---

## Quick Test Script

```bash
#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

SESSION_ID="quick_test_$(date +%s)"

echo "🧪 Testing AI Memory with session: $SESSION_ID"
echo ""

# Test 1
echo "Test 1: Send info"
RESPONSE=$(curl -s -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"My name is Mike and I love Python\",
    \"session_id\": \"$SESSION_ID\",
    \"use_history\": true
  }")

HISTORY_USED=$(echo $RESPONSE | jq -r '.history_used')
if [ "$HISTORY_USED" = "false" ]; then
  echo -e "${GREEN}✅ Test 1 PASSED${NC}"
else
  echo -e "${RED}❌ Test 1 FAILED${NC}"
fi

sleep 1

# Test 2
echo "Test 2: Ask name"
RESPONSE=$(curl -s -X POST http://localhost:3001/api/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"What is my name?\",
    \"session_id\": \"$SESSION_ID\",
    \"use_history\": true
  }")

HISTORY_USED=$(echo $RESPONSE | jq -r '.history_used')
REPLY=$(echo $RESPONSE | jq -r '.reply')

if [ "$HISTORY_USED" = "true" ] && echo "$REPLY" | grep -qi "Mike"; then
  echo -e "${GREEN}✅ Test 2 PASSED${NC}"
  echo "   Reply: $REPLY"
else
  echo -e "${RED}❌ Test 2 FAILED${NC}"
  echo "   Reply: $REPLY"
fi

echo ""
echo "✅ All tests completed!"
```

บันทึกเป็น `quick_test.sh` และรันด้วย:
```bash
chmod +x quick_test.sh
./quick_test.sh
```
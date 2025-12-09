# Redis Caching Strategy for STOU Smart Tour Backend

## Executive Summary

ระบบ STOU Smart Tour มีการเรียก External APIs หลายตัว (Google Search, Places, YouTube, Translate, OpenAI) โดยปัจจุบัน **ไม่มีการ cache ใดๆ** ทำให้:
- ต้นทุน API สูงมาก
- Response time ช้า
- เสี่ยงโดน Rate Limit

**ประมาณการประหยัด**: ลด API calls 50-70% = ประหยัด **$100-150/เดือน**

---

## สถานะปัจจุบัน

### มีพร้อมแล้ว
| Component | File | Status |
|-----------|------|--------|
| Redis Client | `infrastructure/redis/redis.go` | ✅ พร้อมใช้ |
| Cache Keys | `infrastructure/cache/cache_keys.go` | ✅ กำหนดไว้แล้ว |
| TTL Config | `infrastructure/cache/cache_keys.go` | ✅ กำหนดไว้แล้ว |

### ยังไม่ได้ทำ
| Component | Status |
|-----------|--------|
| Search Service Caching | ❌ ไม่มี |
| AI Service Caching | ❌ ไม่มี |
| Utility Service Caching | ❌ ไม่มี |
| Place Details Caching | ❌ ไม่มี |

---

## Cache Keys และ TTL ที่กำหนดไว้

```go
// Prefixes
PrefixSearch       = "search"        // ค้นหาทั่วไป
PrefixSearchAI     = "search:ai"     // AI Search
PrefixPlace        = "place"         // Place ข้อมูลพื้นฐาน
PrefixPlaceDetails = "place:details" // Place รายละเอียด
PrefixNearbyPlaces = "places:nearby" // สถานที่ใกล้เคียง
PrefixYouTube      = "youtube"       // YouTube
PrefixTranslate    = "translate"     // แปลภาษา

// TTLs
TTLSearch       = 1 * time.Hour      // 1 ชั่วโมง
TTLSearchAI     = 6 * time.Hour      // 6 ชั่วโมง
TTLPlace        = 1 * time.Hour      // 1 ชั่วโมง
TTLPlaceDetails = 24 * time.Hour     // 24 ชั่วโมง
TTLNearbyPlaces = 1 * time.Hour      // 1 ชั่วโมง
TTLYouTube      = 6 * time.Hour      // 6 ชั่วโมง
TTLTranslate    = 7 * 24 * time.Hour // 7 วัน
```

---

## จุดที่ต้องเพิ่ม Cache (15 จุด)

### Priority: CRITICAL (ทำก่อน)

#### 1. Place Details - `GetPlaceDetails()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:425` |
| **API** | Google Places Details |
| **Cache Key** | `place:details:{placeID}` |
| **TTL** | 24 ชั่วโมง |
| **Hit Rate** | 40-60% |
| **ประหยัด** | $30-50/เดือน |

#### 2. Nearby Places - `SearchPlaces()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:359` |
| **API** | Google Places Nearby |
| **Cache Key** | `places:nearby:{hash(lat+lng+radius+type+keyword)}` |
| **TTL** | 1 ชั่วโมง |
| **Hit Rate** | 20-30% |
| **ประหยัด** | $18-30/เดือน |

#### 3. Nearby Places - `SearchNearbyPlaces()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:489` |
| **API** | Google Places Nearby |
| **Cache Key** | `places:nearby:{hash(lat+lng+radius+type+keyword)}` |
| **TTL** | 1 ชั่วโมง |
| **Hit Rate** | 20-30% |
| **ประหยัด** | $18-30/เดือน |

#### 4. AI Search - `AISearch()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/ai_service_impl.go:45` |
| **API** | Google Search + OpenAI |
| **Cache Key** | `search:ai:{hash(query)}` |
| **TTL** | 6 ชั่วโมง |
| **Hit Rate** | 50-70% |
| **ประหยัด** | $35-55/เดือน |

#### 5. Translation - `Translate()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/utility_service_impl.go:36` |
| **API** | Google Translate |
| **Cache Key** | `translate:{hash(text+sourceLang+targetLang)}` |
| **TTL** | 7 วัน |
| **Hit Rate** | 40-60% |
| **ประหยัด** | $20-40/เดือน |

---

### Priority: HIGH (ทำหลังจาก CRITICAL)

#### 6. Website Search - `SearchWebsites()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:143` |
| **API** | Google Custom Search |
| **Cache Key** | `search:{hash(query+type)}:{page}` |
| **TTL** | 1 ชั่วโมง |
| **ประหยัด** | $2-3/เดือน |

#### 7. Image Search - `SearchImages()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:178` |
| **API** | Google Custom Search (Images) |
| **Cache Key** | `search:{hash(query+image)}:{page}` |
| **TTL** | 6 ชั่วโมง |
| **ประหยัด** | $1-2/เดือน |

#### 8. Video Search - `SearchVideos()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:226` |
| **API** | YouTube Data API |
| **Cache Key** | `youtube:{hash(query+limit)}` |
| **TTL** | 6 ชั่วโมง |
| **ประหยัด** | $2-3/เดือน |

#### 9. Video Details - `GetVideoDetails()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/search_service_impl.go:309` |
| **API** | YouTube Data API |
| **Cache Key** | `youtube:{videoID}` |
| **TTL** | 24 ชั่วโมง |
| **ประหยัด** | Quota reduction |

#### 10. Create Chat Session - `CreateChatSession()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/ai_service_impl.go:98` |
| **API** | Google Search + OpenAI |
| **Cache Key** | `search:ai:{hash(query)}` (reuse AISearch cache) |
| **TTL** | 6 ชั่วโมง |
| **ประหยัด** | $20-30/เดือน |

---

### Priority: MEDIUM (Optional)

#### 11. Language Detection - `DetectLanguage()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/utility_service_impl.go:63` |
| **API** | Google Translate (Detect) |
| **Cache Key** | `detect:{hash(text)}` |
| **TTL** | 7 วัน |

#### 12. Send Message Search - `SendMessage()`
| รายละเอียด | ค่า |
|------------|-----|
| **File** | `application/serviceimpl/ai_service_impl.go:250` |
| **API** | Google Search (for context) |
| **Cache Key** | `search:msg:{hash(message)}` |
| **TTL** | 2 ชั่วโมง |

#### 13-15. User Data Caching (Folders, Favorites)
| รายละเอียด | ค่า |
|------------|-----|
| **Files** | `folder_service_impl.go`, `favorite_service_impl.go` |
| **API** | Database (ไม่ใช่ External API) |
| **Cache Key** | `user:{userID}:folders`, `user:{userID}:favorites` |
| **TTL** | 5 นาที |
| **ประโยชน์** | ลด DB Load |

---

## ต้นทุน API รายเดือน

### ก่อน Cache
| API | ต้นทุน/เดือน |
|-----|-------------|
| Google Custom Search | $15-20 |
| Google Places | $60-80 |
| YouTube Data | $5-10 |
| Google Translate | $40-60 |
| OpenAI GPT-4o-mini | $100-150 |
| **รวม** | **$220-320** |

### หลัง Cache
| API | ต้นทุน/เดือน | ประหยัด |
|-----|-------------|---------|
| Google Custom Search | $9-12 | 30-40% |
| Google Places | $35-50 | 40-50% |
| YouTube Data | $3-6 | 30-40% |
| Google Translate | $25-35 | 40-50% |
| OpenAI GPT-4o-mini | $40-60 | 50-70% |
| **รวม** | **$112-163** | **50-65%** |

### สรุป
- **ประหยัด**: $108-157/เดือน
- **ประหยัดต่อปี**: $1,200-1,800

---

## Implementation Pattern

### ตัวอย่าง: SearchWebsites with Caching

```go
func (s *SearchServiceImpl) SearchWebsites(ctx context.Context, userID uuid.UUID, req *dto.SearchRequest) (*dto.WebsiteSearchResponse, error) {
    // 1. สร้าง Cache Key
    cacheKey := cache.SearchKey(req.Query, "website", req.Page)

    // 2. ตรวจสอบ Cache ก่อน
    var cachedResult dto.WebsiteSearchResponse
    if err := s.redisClient.Get(ctx, cacheKey, &cachedResult); err == nil {
        // Cache HIT - return cached data
        return &cachedResult, nil
    }

    // 3. Cache MISS - เรียก API
    searchResponse, err := s.googleSearch.SearchAll(ctx, req.Query, req.Page, req.PageSize)
    if err != nil {
        return nil, err
    }

    // 4. สร้าง Response
    var websiteResults []dto.WebsiteResult
    for _, r := range searchResponse.Items {
        websiteResults = append(websiteResults, dto.WebsiteResult{
            Title:       r.Title,
            URL:         r.Link,
            Snippet:     r.Snippet,
            DisplayLink: r.DisplayLink,
        })
    }

    response := &dto.WebsiteSearchResponse{
        Query:      req.Query,
        Results:    websiteResults,
        TotalCount: int64(len(websiteResults)),
        Page:       req.Page,
        PageSize:   req.PageSize,
    }

    // 5. บันทึกลง Cache
    _ = s.redisClient.Set(ctx, cacheKey, response, cache.TTLSearch)

    // 6. บันทึก Search History (ไม่เกี่ยวกับ cache)
    s.saveSearchHistory(ctx, userID, req.Query, models.SearchTypeWebsite, len(websiteResults))

    return response, nil
}
```

---

## Cache Helper Functions ที่ต้องเพิ่ม

เพิ่มใน `infrastructure/cache/cache_keys.go`:

```go
// DetectKey generates cache key for language detection
func DetectKey(text string) string {
    return fmt.Sprintf("detect:%s", hashString(text))
}

// VideoDetailsKey generates cache key for video details
func VideoDetailsKey(videoID string) string {
    return fmt.Sprintf("%s:details:%s", PrefixYouTube, videoID)
}

// ChatSearchKey generates cache key for chat message search context
func ChatSearchKey(message string) string {
    return fmt.Sprintf("search:msg:%s", hashString(message))
}

// UserFoldersKey generates cache key for user's folder list
func UserFoldersKey(userID string) string {
    return fmt.Sprintf("user:%s:folders", userID)
}

// UserFavoritesKey generates cache key for user's favorites
func UserFavoritesKey(userID, favoriteType string) string {
    if favoriteType == "" {
        return fmt.Sprintf("user:%s:favorites", userID)
    }
    return fmt.Sprintf("user:%s:favorites:%s", userID, favoriteType)
}
```

---

## Redis Memory Estimation

| Cache Type | Avg Size | Daily Entries | Memory |
|------------|----------|---------------|--------|
| Search Results | 15-20 KB | 100-200 | 2-4 MB |
| Place Details | 30-50 KB | 100-200 | 3-10 MB |
| Nearby Places | 20-30 KB | 60-100 | 1.2-3 MB |
| AI Responses | 10-15 KB | 50-100 | 0.5-1.5 MB |
| Translations | 1-3 KB | 50-100 | 0.05-0.3 MB |
| Videos | 12-18 KB | 100-200 | 1.2-3.6 MB |
| **Total** | | | **~10-25 MB** |

**Recommended Redis maxmemory**: 256-512 MB (ให้มี buffer)

---

## Implementation Order

### Sprint 1 (CRITICAL) - ประหยัด ~$100/เดือน
1. ✅ Place Details Caching
2. ✅ Nearby Places Caching (ทั้ง 2 functions)
3. ✅ AI Search Caching
4. ✅ Translation Caching

### Sprint 2 (HIGH) - ประหยัดเพิ่ม ~$30/เดือน
5. Website Search Caching
6. Image Search Caching
7. Video Search Caching
8. Video Details Caching
9. Create Chat Session Caching

### Sprint 3 (MEDIUM) - ลด Load
10. Language Detection Caching
11. Chat Message Search Caching
12. User Folders Caching
13. User Favorites Caching

---

## Error Handling Pattern

```go
// ถ้า Redis ล่ม ให้ fallback ไป API ปกติ
var cachedResult dto.SomeResponse
err := s.redisClient.Get(ctx, cacheKey, &cachedResult)
if err == nil {
    return &cachedResult, nil
}

// Log cache miss/error แต่ไม่ block
if err != redis.Nil {
    log.Printf("Cache error for key %s: %v", cacheKey, err)
}

// Continue to API call...
```

---

## Monitoring

### Metrics ที่ควรเก็บ
- Cache Hit Rate (ควรได้ 40-60%)
- Cache Miss Rate
- Redis Memory Usage
- API Calls per hour (before/after cache)

### Redis Commands สำหรับ Monitor
```bash
# ดู memory usage
redis-cli INFO memory

# ดู key count
redis-cli DBSIZE

# ดู hit/miss rate
redis-cli INFO stats | grep keyspace
```

---

## Files ที่ต้องแก้ไข

| File | Changes |
|------|---------|
| `infrastructure/cache/cache_keys.go` | เพิ่ม helper functions |
| `application/serviceimpl/search_service_impl.go` | เพิ่ม cache logic ทุก function |
| `application/serviceimpl/ai_service_impl.go` | เพิ่ม cache logic |
| `application/serviceimpl/utility_service_impl.go` | เพิ่ม cache logic |
| `application/serviceimpl/folder_service_impl.go` | เพิ่ม cache + invalidation |
| `application/serviceimpl/favorite_service_impl.go` | เพิ่ม cache + invalidation |

---

## Conclusion

การ implement Redis caching จะ:
1. **ลดต้นทุน API** $100-150/เดือน
2. **เพิ่มความเร็ว** 50-100ms สำหรับ cached queries
3. **ลดความเสี่ยง Rate Limit** 50-65% fewer API calls
4. **รองรับ Users มากขึ้น** 2-3x concurrent users ด้วย quota เดิม

**ROI**: ใช้เวลาพัฒนา 8-12 ชั่วโมง, คืนทุนภายใน 1 สัปดาห์

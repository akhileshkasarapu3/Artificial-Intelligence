from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter(prefix="/tools")

class CalendarSlotsRequest(BaseModel):
    date: str  # YYYY-MM-DD

@router.post("/get_calendar_slots")
def get_calendar_slots(req: CalendarSlotsRequest):
    # Mock data for Day 1
    return {"date": req.date, "slots": ["10:00", "11:30", "15:00"]}

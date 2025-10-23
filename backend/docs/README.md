# 📚 Backend Documentation

This directory contains comprehensive documentation for the Python FastAPI backend.

## 📖 Available Guides

### **[Complete Setup Guide](COMPLETE_SETUP_GUIDE.md)**
- Step-by-step setup instructions
- Environment configuration
- Service startup and verification
- Testing procedures
- API documentation links

### **[Troubleshooting Guide](TROUBLESHOOTING.md)**
- Common issues and solutions
- Docker troubleshooting
- Database connection problems
- Service startup failures
- API error debugging

### **[Project Vision & Roadmap](PROJECT_VISION.md)**
- Project vision and goals
- Feature roadmap and ideas
- Future development plans

## 🚀 Quick Reference

**Start services:**
```bash
make dev
```

**Run tests:**
```bash
./test-setup.sh
```

**Check status:**
```bash
make status
```

**View logs:**
```bash
make logs
```

## 🔗 API Documentation

Once services are running:
- **User Service**: http://localhost:8001/docs
- **Plaid Service**: http://localhost:8002/docs  
- **Transaction Service**: http://localhost:8003/docs

#!/usr/bin/env python3
"""
Basic structure test for AI Context Gap Tracker MCP Server
Tests the package structure and configuration without requiring external dependencies
"""

import os
import sys
import json

def test_file_structure():
    """Test if all required files exist"""
    print("🔍 Testing file structure...")
    
    required_files = [
        "pyproject.toml",
        "mcp_server/__init__.py",
        "mcp_server/main.py",
        "mcp-server/requirements.txt",
        "mcp-server/Dockerfile",
        "docs/MCP_SERVER_SETUP.md",
        "docker-compose.yml"
    ]
    
    missing_files = []
    for file_path in required_files:
        if os.path.exists(file_path):
            print(f"✅ {file_path}")
        else:
            print(f"❌ {file_path}")
            missing_files.append(file_path)
    
    return len(missing_files) == 0

def test_pyproject_config():
    """Test pyproject.toml configuration"""
    print("\n🔍 Testing pyproject.toml configuration...")
    
    try:
        with open("pyproject.toml", "r") as f:
            content = f.read()
        
        required_sections = [
            "[project]",
            "[project.scripts]",
            "ai-context-gap-tracker",
            "mcp_server.main:main"
        ]
        
        for section in required_sections:
            if section in content:
                print(f"✅ Found: {section}")
            else:
                print(f"❌ Missing: {section}")
                return False
                
        print("✅ pyproject.toml is properly configured")
        return True
        
    except FileNotFoundError:
        print("❌ pyproject.toml not found")
        return False
    except Exception as e:
        print(f"❌ pyproject.toml test error: {e}")
        return False

def test_docker_compose():
    """Test docker-compose.yml includes MCP server"""
    print("\n🔍 Testing docker-compose.yml...")
    
    try:
        with open("docker-compose.yml", "r") as f:
            content = f.read()
        
        if "mcp-server:" in content:
            print("✅ MCP server service defined")
        else:
            print("❌ MCP server service not found")
            return False
            
        if "TRACKER_API_ENDPOINT" in content:
            print("✅ Environment variables configured")
        else:
            print("❌ Environment variables missing")
            return False
            
        return True
        
    except FileNotFoundError:
        print("❌ docker-compose.yml not found")
        return False
    except Exception as e:
        print(f"❌ docker-compose.yml test error: {e}")
        return False

def test_mcp_server_structure():
    """Test MCP server Python files structure"""
    print("\n🔍 Testing MCP server structure...")
    
    try:
        # Test __init__.py
        with open("mcp_server/__init__.py", "r") as f:
            init_content = f.read()
        
        if "MCPServer" in init_content:
            print("✅ __init__.py exports MCPServer")
        else:
            print("❌ __init__.py doesn't export MCPServer")
            return False
        
        # Test main.py
        with open("mcp_server/main.py", "r") as f:
            main_content = f.read()
        
        required_in_main = [
            "class MCPServer",
            "def main()",
            "rewrite_prompt",
            "audit_response",
            "track_context"
        ]
        
        for item in required_in_main:
            if item in main_content:
                print(f"✅ Found in main.py: {item}")
            else:
                print(f"❌ Missing in main.py: {item}")
                return False
        
        return True
        
    except FileNotFoundError as e:
        print(f"❌ File not found: {e}")
        return False
    except Exception as e:
        print(f"❌ Structure test error: {e}")
        return False

def test_documentation():
    """Test if documentation exists and contains key sections"""
    print("\n🔍 Testing documentation...")
    
    try:
        with open("docs/MCP_SERVER_SETUP.md", "r") as f:
            doc_content = f.read()
        
        required_sections = [
            "uvx",
            "Claude Desktop",
            "claude_desktop_config.json",
            "TRACKER_API_ENDPOINT"
        ]
        
        for section in required_sections:
            if section in doc_content:
                print(f"✅ Documentation includes: {section}")
            else:
                print(f"❌ Documentation missing: {section}")
                return False
        
        return True
        
    except FileNotFoundError:
        print("❌ MCP_SERVER_SETUP.md not found")
        return False
    except Exception as e:
        print(f"❌ Documentation test error: {e}")
        return False

def main():
    """Run all structure tests"""
    print("🚀 AI Context Gap Tracker MCP Server - Structure Tests")
    print("=" * 60)
    
    tests = [
        ("File Structure", test_file_structure),
        ("pyproject.toml Configuration", test_pyproject_config),
        ("Docker Compose Configuration", test_docker_compose),
        ("MCP Server Code Structure", test_mcp_server_structure),
        ("Documentation", test_documentation),
    ]
    
    passed = 0
    total = len(tests)
    
    for test_name, test_func in tests:
        try:
            result = test_func()
            if result:
                passed += 1
        except Exception as e:
            print(f"❌ Test '{test_name}' failed with exception: {e}")
    
    print("\n" + "=" * 60)
    print(f"📊 Test Results: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 All structure tests passed! MCP server structure is correct.")
        print("\n📋 Next steps to test full functionality:")
        print("1. Install dependencies: pip install httpx")
        print("2. Start backend services: docker-compose up -d")
        print("3. Test with: python3 scripts/test_mcp_server.py")
        print("4. Or run directly: uvx --from . ai-context-gap-tracker")
    else:
        print("⚠️  Some structure tests failed. Please fix the issues above.")
    
    return passed == total

if __name__ == "__main__":
    success = main()
    sys.exit(0 if success else 1)
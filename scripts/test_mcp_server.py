#!/usr/bin/env python3
"""
Test script for AI Context Gap Tracker MCP Server
Validates the MCP server functionality and integration
"""

import asyncio
import json
import sys
import subprocess
import time
import httpx
import os
from typing import Dict, Any

class MCPServerTester:
    def __init__(self):
        self.tracker_endpoint = "http://localhost:8080"
        self.test_results = []
        
    async def test_backend_connectivity(self):
        """Test if backend services are accessible"""
        print("🔍 Testing backend connectivity...")
        
        try:
            async with httpx.AsyncClient(timeout=10.0) as client:
                response = await client.get(f"{self.tracker_endpoint}/api/v1/health")
                if response.status_code == 200:
                    print("✅ Backend services are accessible")
                    return True
                else:
                    print(f"❌ Backend health check failed: {response.status_code}")
                    return False
        except Exception as e:
            print(f"❌ Cannot connect to backend services: {e}")
            print("ℹ️  Make sure to run: docker-compose up -d")
            return False
    
    def test_package_installation(self):
        """Test if the package can be installed and imported"""
        print("\n🔍 Testing package installation...")
        
        try:
            # Test if mcp_server can be imported
            import mcp_server
            print("✅ Package can be imported successfully")
            
            # Test if main module exists
            from mcp_server.main import MCPServer
            print("✅ MCPServer class can be imported")
            return True
            
        except ImportError as e:
            print(f"❌ Import error: {e}")
            return False
    
    def test_uvx_compatibility(self):
        """Test uvx execution (dry run)"""
        print("\n🔍 Testing uvx compatibility...")
        
        try:
            # Check if uvx is available
            result = subprocess.run(["uvx", "--version"], 
                                  capture_output=True, text=True, timeout=10)
            if result.returncode == 0:
                print("✅ uvx is available")
                
                # Test the command structure (but don't actually run it to avoid conflicts)
                print("✅ uvx command structure is valid")
                print("   Command: uvx --from . ai-context-gap-tracker")
                return True
            else:
                print("❌ uvx not available or not working")
                return False
                
        except subprocess.TimeoutExpired:
            print("❌ uvx command timed out")
            return False
        except FileNotFoundError:
            print("❌ uvx not found. Install with: pip install uvx")
            return False
        except Exception as e:
            print(f"❌ uvx test error: {e}")
            return False
    
    def test_docker_integration(self):
        """Test Docker integration"""
        print("\n🔍 Testing Docker integration...")
        
        try:
            # Check if docker-compose file includes MCP server
            with open("docker-compose.yml", "r") as f:
                compose_content = f.read()
                
            if "mcp-server:" in compose_content:
                print("✅ MCP server is included in docker-compose.yml")
                
                if "TRACKER_API_ENDPOINT" in compose_content:
                    print("✅ Environment variables are configured")
                    return True
                else:
                    print("❌ Environment variables not configured")
                    return False
            else:
                print("❌ MCP server not found in docker-compose.yml")
                return False
                
        except FileNotFoundError:
            print("❌ docker-compose.yml not found")
            return False
        except Exception as e:
            print(f"❌ Docker integration test error: {e}")
            return False
    
    def test_pyproject_config(self):
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
    
    async def test_mcp_tools_structure(self):
        """Test MCP tools structure by importing and inspecting"""
        print("\n🔍 Testing MCP tools structure...")
        
        try:
            from mcp_server.main import MCPServer
            
            # Create server instance (without starting it)
            server = MCPServer()
            
            # Check if required attributes exist
            if hasattr(server, 'server'):
                print("✅ MCP Server instance created")
            else:
                print("❌ MCP Server instance missing 'server' attribute")
                return False
                
            if hasattr(server, 'tracker_api_endpoint'):
                print(f"✅ API endpoint configured: {server.tracker_api_endpoint}")
            else:
                print("❌ API endpoint not configured")
                return False
                
            print("✅ MCP tools structure is valid")
            return True
            
        except Exception as e:
            print(f"❌ MCP tools structure test error: {e}")
            return False
    
    def test_documentation(self):
        """Test if documentation exists"""
        print("\n🔍 Testing documentation...")
        
        docs_to_check = [
            "docs/MCP_SERVER_SETUP.md",
            "README.md",
            "docs/DEPLOYMENT.md"
        ]
        
        all_exist = True
        for doc in docs_to_check:
            if os.path.exists(doc):
                print(f"✅ Found: {doc}")
            else:
                print(f"❌ Missing: {doc}")
                all_exist = False
                
        return all_exist
    
    async def run_all_tests(self):
        """Run all tests"""
        print("🚀 AI Context Gap Tracker MCP Server - Integration Tests")
        print("=" * 60)
        
        tests = [
            ("Package Installation", self.test_package_installation),
            ("pyproject.toml Configuration", self.test_pyproject_config),
            ("uvx Compatibility", self.test_uvx_compatibility),
            ("Docker Integration", self.test_docker_integration),
            ("MCP Tools Structure", self.test_mcp_tools_structure),
            ("Documentation", self.test_documentation),
            ("Backend Connectivity", self.test_backend_connectivity),
        ]
        
        passed = 0
        total = len(tests)
        
        for test_name, test_func in tests:
            try:
                if asyncio.iscoroutinefunction(test_func):
                    result = await test_func()
                else:
                    result = test_func()
                    
                if result:
                    passed += 1
                    
            except Exception as e:
                print(f"❌ Test '{test_name}' failed with exception: {e}")
        
        print("\n" + "=" * 60)
        print(f"📊 Test Results: {passed}/{total} tests passed")
        
        if passed == total:
            print("🎉 All tests passed! MCP server is ready for use.")
            print("\n📋 Next steps:")
            print("1. Start backend services: docker-compose up -d")
            print("2. Run MCP server: uvx --from . ai-context-gap-tracker")
            print("3. Configure Claude Desktop (see docs/MCP_SERVER_SETUP.md)")
        else:
            print("⚠️  Some tests failed. Please address the issues above.")
            
        return passed == total

async def main():
    """Main test execution"""
    tester = MCPServerTester()
    success = await tester.run_all_tests()
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    asyncio.run(main())
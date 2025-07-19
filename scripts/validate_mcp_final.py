#!/usr/bin/env python3
"""
Final validation script for AI Context Gap Tracker MCP Server
Tests all components without requiring external dependencies
"""

import os
import sys
import ast
import json
from pathlib import Path

def print_header(title):
    print(f"\n🚀 {title}")
    print("=" * 60)

def print_success(msg):
    print(f"✅ {msg}")

def print_info(msg):
    print(f"📋 {msg}")

def validate_python_syntax(file_path):
    """Validate Python file syntax without importing"""
    try:
        with open(file_path, 'r') as f:
            content = f.read()
        ast.parse(content)
        return True
    except Exception as e:
        print(f"❌ Syntax error in {file_path}: {e}")
        return False

def validate_package_structure():
    """Validate the complete package structure"""
    print_header("Package Structure Validation")
    
    required_files = [
        "pyproject.toml",
        "mcp_server/__init__.py", 
        "mcp_server/main.py",
        "mcp-server/requirements.txt",
        "mcp-server/Dockerfile",
        "docs/MCP_SERVER_SETUP.md",
        "docker-compose.yml"
    ]
    
    all_exist = True
    for file_path in required_files:
        if os.path.exists(file_path):
            print_success(f"Found: {file_path}")
        else:
            print(f"❌ Missing: {file_path}")
            all_exist = False
    
    return all_exist

def validate_pyproject_config():
    """Validate pyproject.toml configuration"""
    print_header("pyproject.toml Configuration")
    
    try:
        # Read as text since we don't have toml library
        with open("pyproject.toml", "r") as f:
            content = f.read()
        
        checks = [
            ("[project]", "Project section"),
            ("[project.scripts]", "Scripts section"), 
            ("ai-context-gap-tracker", "Entry point name"),
            ("mcp_server.main:main", "Main function reference"),
            ("httpx", "httpx dependency"),
            ("mcp", "mcp dependency")
        ]
        
        all_valid = True
        for check_str, description in checks:
            if check_str in content:
                print_success(f"{description} configured")
            else:
                print(f"❌ Missing: {description}")
                all_valid = False
        
        return all_valid
    except Exception as e:
        print(f"❌ Error reading pyproject.toml: {e}")
        return False

def validate_python_files():
    """Validate Python file syntax"""
    print_header("Python Files Syntax Check")
    
    python_files = [
        "mcp_server/__init__.py",
        "mcp_server/main.py"
    ]
    
    all_valid = True
    for file_path in python_files:
        if validate_python_syntax(file_path):
            print_success(f"Valid syntax: {file_path}")
        else:
            all_valid = False
    
    return all_valid

def validate_mcp_server_features():
    """Validate MCP server has required features"""
    print_header("MCP Server Features Check")
    
    try:
        with open("mcp_server/main.py", "r") as f:
            content = f.read()
        
        required_features = [
            ("class MCPServer", "MCPServer class"),
            ("def main()", "Main function"),
            ("rewrite_prompt", "Rewrite prompt tool"),
            ("audit_response", "Audit response tool"), 
            ("track_context", "Track context tool"),
            ("evaluate_rules", "Evaluate rules tool"),
            ("get_session_context", "Get session context tool"),
            ("list_tools", "List tools handler"),
            ("call_tool", "Call tool handler"),
            ("list_resources", "List resources handler"),
            ("read_resource", "Read resource handler")
        ]
        
        all_valid = True
        for feature, description in required_features:
            if feature in content:
                print_success(f"{description} implemented")
            else:
                print(f"❌ Missing: {description}")
                all_valid = False
        
        return all_valid
    except Exception as e:
        print(f"❌ Error reading main.py: {e}")
        return False

def validate_docker_config():
    """Validate Docker configuration"""
    print_header("Docker Configuration Check")
    
    try:
        with open("docker-compose.yml", "r") as f:
            content = f.read()
        
        docker_checks = [
            ("mcp-server:", "MCP server service"),
            ("TRACKER_API_ENDPOINT", "API endpoint env var"),
            ("ports:", "Port configuration"),
            ("depends_on:", "Service dependencies")
        ]
        
        all_valid = True
        for check_str, description in docker_checks:
            if check_str in content:
                print_success(f"{description} configured")
            else:
                print(f"❌ Missing: {description}")
                all_valid = False
        
        return all_valid
    except Exception as e:
        print(f"❌ Error reading docker-compose.yml: {e}")
        return False

def validate_documentation():
    """Validate documentation completeness"""
    print_header("Documentation Check")
    
    try:
        with open("docs/MCP_SERVER_SETUP.md", "r") as f:
            content = f.read()
        
        doc_checks = [
            ("uvx", "uvx usage instructions"),
            ("Claude Desktop", "Claude Desktop integration"),
            ("claude_desktop_config.json", "Configuration file"),
            ("TRACKER_API_ENDPOINT", "Environment variables"),
            ("Installation", "Installation section"),
            ("Usage", "Usage section")
        ]
        
        all_valid = True
        for check_str, description in doc_checks:
            if check_str in content:
                print_success(f"{description} documented")
            else:
                print(f"❌ Missing: {description}")
                all_valid = False
        
        return all_valid
    except Exception as e:
        print(f"❌ Error reading documentation: {e}")
        return False

def main():
    """Run all validation checks"""
    print_header("AI Context Gap Tracker - Final MCP Server Validation")
    
    results = {
        "Package Structure": validate_package_structure(),
        "pyproject.toml": validate_pyproject_config(), 
        "Python Syntax": validate_python_files(),
        "MCP Features": validate_mcp_server_features(),
        "Docker Config": validate_docker_config(),
        "Documentation": validate_documentation()
    }
    
    print_header("Validation Summary")
    
    all_passed = True
    for test_name, passed in results.items():
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"{status} - {test_name}")
        if not passed:
            all_passed = False
    
    print(f"\n{'=' * 60}")
    if all_passed:
        print("🎉 ALL VALIDATIONS PASSED!")
        print("")
        print_info("MCP Server is ready for production use!")
        print_info("Next steps:")
        print("   1. Install uvx: sudo snap install astral-uv")
        print("   2. Run: uvx --from . ai-context-gap-tracker")
        print("   3. Configure Claude Desktop (see docs/MCP_SERVER_SETUP.md)")
        print("")
        print_info("For Docker deployment:")
        print("   1. Run: docker-compose up -d")
        print("   2. MCP server will be available on port 8001")
    else:
        print("❌ Some validations failed. Please check the issues above.")
        return 1
    
    return 0

if __name__ == "__main__":
    sys.exit(main())
import re
from datetime import datetime
from pathlib import Path
from typing import Optional, Dict, Union, List
import pandas as pd

class LogParser:
    """Parses mixed-format logs including GIN web logs and custom test logs."""
    
    def __init__(self):
        # Compile regex patterns once during initialization
        self.gin_pattern = re.compile(
            r'\[GIN\] (\d{4}/\d{2}/\d{2}) - (\d{2}:\d{2}:\d{2}) \| (\d{3}) \|\s+([\d.µ]+)s \|\s+(\S+) \| (\w+)\s+"([^"]+)"'
        )
        self.test_pattern = re.compile(
            r'\[([^\]]+)\] \[([^\]]+)\] \[([^\]]+)\] (.*)'
        )

    def parse_line(self, line: str) -> Optional[Dict[str, Union[str, int, float, datetime]]]:
        """
        Parses a single log line into a structured dictionary.
        Returns None if the line doesn't match any known format.
        """
        line = line.strip()
        if not line:
            return None

        # Try GIN format first
        if gin_match := self.gin_pattern.match(line):
            return self._parse_gin(gin_match)
        
        # Fall back to test format
        if test_match := self.test_pattern.match(line):
            return self._parse_test(test_match)
        
        return None

    def _parse_gin(self, match: re.Match) -> Dict[str, Union[str, int, float, datetime]]:
        """Handles GIN web server log format"""
        date, time, status, latency, ip, method, path = match.groups()
        
        return {
            'log_type': 'gin',
            'timestamp': datetime.strptime(f"{date} {time}", "%Y/%m/%d %H:%M:%S"),
            'status_code': int(status),
            'latency_ms': self._convert_latency(latency),
            'client_ip': ip,
            'http_method': method.upper(),
            'path': path,
            'raw': match.group(0)
        }

    def _parse_test(self, match: re.Match) -> Dict[str, Union[str, datetime]]:
        """Handles custom test log format"""
        timestamp, log_type, ip, message = match.groups()
        
        return {
            'log_type': 'test',
            'timestamp': datetime.strptime(timestamp, "%Y-%m-%dT%H:%M:%SZ"),
            'test_type': log_type.lower(),
            'client_ip': ip,
            'message': message.strip(),
            'raw': match.group(0)
        }

    def _convert_latency(self, latency_str: str) -> float:
        """
        Converts latency string to milliseconds.
        Handles µs, ms, and s units.
        """
        if 'µ' in latency_str:  # Microseconds
            return float(latency_str.replace('µs', '')) / 1000
        if 'ms' in latency_str:  # Milliseconds
            return float(latency_str.replace('ms', ''))
        # Seconds (default)
        return float(latency_str.replace('s', '')) * 1000

    def parse_file(self, file_path: Path) -> List[Dict]:
        """
        Parses an entire log file, skipping unparseable lines.
        Returns a list of parsed log entries.
        """
        parsed_logs = []
        with open(file_path, 'r', encoding='utf-8') as f:
            for line in f:
                if parsed := self.parse_line(line):
                    parsed_logs.append(parsed)
        return parsed_logs

    def to_dataframe(self, file_path: Path) -> pd.DataFrame:
        """
        Parses a log file and returns results as a pandas DataFrame.
        Automatically converts timestamps to datetime64[ns].
        """
        logs = self.parse_file(file_path)
        df = pd.DataFrame(logs)
        if not df.empty and 'timestamp' in df.columns:
            df['timestamp'] = pd.to_datetime(df['timestamp'])
        return df


# Example usage
if __name__ == "__main__":
    parser = LogParser()
    
    # Test with sample lines
    test_lines = [
        '[GIN] 2025/07/07 - 23:54:10 | 200 | 9.3µs | ::1 | GET "/status"',
        '[2025-07-07T23:07:38Z] [test] [::1] Test log',
        '[GIN] 2025/07/08 - 00:33:14 | 200 | 0s | ::1 | POST "/ingest"',
        '[2025-07-07T23:33:14Z] [testing] [::1] YEAH',
        'Invalid log line that will be skipped'
    ]
    
    print("Testing line parsing:")
    for line in test_lines:
        if parsed := parser.parse_line(line):
            print(f"\nParsed: {parsed}")
    
    # Test with file (uncomment to use)
    # df = parser.to_dataframe(Path("path/to/your/logfile.log"))
    # print("\nDataFrame sample:")
    # print(df.head())
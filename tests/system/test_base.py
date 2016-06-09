from logstashbeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Logstashbeat normally
        """
        self.render_config_template(
                path=os.path.abspath(self.working_dir) + "/log/*"
        )

        logstashbeat_proc = self.start_beat()
        self.wait_until( lambda: self.log_contains("logstashbeat is running"))
        exit_code = logstashbeat_proc.kill_and_wait()
        assert exit_code == 0

import { Navbar, NavbarBrand, NavbarContent, NavbarItem, Button, Link, Tooltip } from "@nextui-org/react";
import { FaGithub as GithubIcon } from "react-icons/fa6";
import { MdDns as DnsIcon } from "react-icons/md";
import { useTranslation } from "react-i18next";

import ThemeSwitcher from "./ThemeSwitcher";
import LangSwitcher from "./LangSwitcher";
import Upload from "./Upload";

export const NavBar = () => {
  const { t } = useTranslation();
  return (
    <Navbar isBordered isBlurred shouldHideOnScroll>
      <NavbarBrand>
        <Link href="/" color="foreground">
          <DnsIcon className="w-6 h-6 mr-2" />
          <p className="font-bold text-inherit">{t("title")}</p>
        </Link>
      </NavbarBrand>

      <NavbarContent justify="end">
        <NavbarItem>
          <Tooltip content={t("tip_github")}>
            <Link href="https://github.com/xxnuo/dns-benchmark" target="_blank">
              <Button variant="ghost" aria-label={t("tip_github")}>
                <GithubIcon />
                <span className="ml-2">{t("tip_github")}</span>
              </Button>
            </Link>
          </Tooltip>
        </NavbarItem>
        <NavbarItem>
          <LangSwitcher />
        </NavbarItem>
        <NavbarItem>
          <ThemeSwitcher />
        </NavbarItem>
        <NavbarItem>
          <Upload />
        </NavbarItem>
      </NavbarContent>
    </Navbar>
  );
};
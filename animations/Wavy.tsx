import { FC } from "react";
import { motion, Variants, HTMLMotionProps } from "framer-motion";


interface Props extends HTMLMotionProps<"div"> {
    text: string;
    delay?: number;
    duration?: number;
}
// const Wavy: FC<HTMLMotionProps> = motion.custom(motion.div);

const Letters: FC<Props> = ({
    text,
    delay = 0,
    duration = 0.06,
    ...props  
}: Props) => {
    const letters = Array.from(text)

    const container: Variants = {
        hidden: {
          opacity: 0,
        },
        show: (i: number = 1) => ({
          opacity: 1,
          transition: { staggerChildren: duration, delayChildren: i * delay },
        }),
      }

    const child: Variants = {
        show: {
            y: 0,
            opacity: 1,
            transition: {
                type: 'spring',
                damping: 10,
            },
        },
        hidden: {
            y: '100%',
            opacity: 1
        }
    }

    return (
        <>
        {heading === 'h1' ? (

        )}
        </>
    )
}